package com.tourtect.feature.assistant

import androidx.lifecycle.ViewModel
import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.viewModelScope
import com.tourtect.core.model.AssistantConfirmation
import com.tourtect.core.model.AssistantItem
import com.tourtect.core.network.AssistantRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import java.util.Locale
import java.util.UUID
import javax.inject.Inject
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.Job
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch

sealed interface RealtimeState {
    data object Disconnected : RealtimeState
    data class Unavailable(val reason: String) : RealtimeState
}

data class UiError(
    val message: String,
    val retryable: Boolean
)

data class AssistantUiState(
    val sessionId: String? = null,
    val messages: List<AssistantItem> = emptyList(),
    val input: String = "",
    val isSending: Boolean = false,
    val realtimeState: RealtimeState = RealtimeState.Unavailable(
        "Live voice is not connected to assistant sessions in this build."
    ),
    val pendingConfirmation: AssistantConfirmation? = null,
    val feedbackSubmitted: Set<String> = emptySet(),
    val error: UiError? = null
)

@HiltViewModel
class AssistantViewModel @Inject constructor(
    private val repository: AssistantRepository,
    private val savedStateHandle: SavedStateHandle
) : ViewModel() {
    private val locale = Locale.getDefault().toLanguageTag()
    private val _uiState = MutableStateFlow(
        AssistantUiState(
            sessionId = savedStateHandle[SESSION_ID_KEY],
            input = savedStateHandle[INPUT_KEY].orEmpty()
        )
    )
    val uiState: StateFlow<AssistantUiState> = _uiState.asStateFlow()

    private var requestJob: Job? = null

    init {
        restoreOrCreateSession()
    }

    fun updateInput(value: String) {
        savedStateHandle[INPUT_KEY] = value
        _uiState.update { it.copy(input = value, error = null) }
    }

    fun usePrompt(prompt: String) {
        savedStateHandle[INPUT_KEY] = prompt
        _uiState.update { it.copy(input = prompt, error = null) }
    }

    fun sendMessage() {
        val text = _uiState.value.input.trim()
        if (text.isEmpty() || _uiState.value.isSending) return

        val messageId = UUID.randomUUID().toString()
        requestJob = viewModelScope.launch {
            _uiState.update {
                it.copy(
                    input = "",
                    isSending = true,
                    error = null,
                    messages = it.messages + AssistantItem.UserMessage(messageId, text)
                )
            }
            savedStateHandle[INPUT_KEY] = ""
            try {
                val sessionId = ensureSession()
                val response = repository.sendText(sessionId, text, locale, messageId)
                _uiState.update {
                    it.copy(
                        isSending = false,
                        messages = it.messages + AssistantItem.Response(response),
                        pendingConfirmation = response.requestedConfirmation
                    )
                }
            } catch (cancelled: CancellationException) {
                _uiState.update { it.copy(isSending = false) }
                throw cancelled
            } catch (_: Exception) {
                savedStateHandle[INPUT_KEY] = text
                _uiState.update {
                    it.copy(
                        input = text,
                        isSending = false,
                        error = UiError(
                            message = "Tourtect cannot interpret this automatically right now.",
                            retryable = true
                        )
                    )
                }
            }
        }
    }

    fun cancelRequest() {
        requestJob?.cancel()
        _uiState.update { it.copy(isSending = false) }
    }

    fun resetSession() {
        requestJob?.cancel()
        val previousSession = _uiState.value.sessionId
        requestJob = viewModelScope.launch {
            savedStateHandle.remove<String>(SESSION_ID_KEY)
            savedStateHandle[INPUT_KEY] = ""
            _uiState.value = AssistantUiState(isSending = true)
            if (previousSession != null) {
                runCatching { repository.deleteSession(previousSession) }
            }
            try {
                ensureSession()
                _uiState.update { it.copy(isSending = false) }
            } catch (_: Exception) {
                showSessionUnavailable()
            }
        }
    }

    fun decideConfirmation(confirmed: Boolean) {
        val sessionId = _uiState.value.sessionId ?: return
        val confirmation = _uiState.value.pendingConfirmation ?: return
        if (_uiState.value.isSending) return

        requestJob = viewModelScope.launch {
            _uiState.update { it.copy(isSending = true, error = null) }
            try {
                val result = repository.confirmAction(
                    sessionId = sessionId,
                    confirmationId = confirmation.confirmationId,
                    confirmed = confirmed
                )
                _uiState.update {
                    it.copy(
                        isSending = false,
                        pendingConfirmation = null,
                        messages = it.messages + AssistantItem.StatusNotice(
                            stableId = "confirmation-${result.confirmationId}",
                            title = "Confirmation recorded",
                            description = "The backend marked ${result.action} as ${result.status}."
                        )
                    )
                }
            } catch (_: Exception) {
                _uiState.update {
                    it.copy(
                        isSending = false,
                        error = UiError("The confirmation could not be verified.", retryable = true)
                    )
                }
            }
        }
    }

    fun submitFeedback(assistantMessageId: String, helpful: Boolean) {
        val sessionId = _uiState.value.sessionId ?: return
        if (assistantMessageId in _uiState.value.feedbackSubmitted) return

        viewModelScope.launch {
            try {
                repository.submitFeedback(
                    sessionId = sessionId,
                    assistantMessageId = assistantMessageId,
                    feedbackType = if (helpful) "helpful" else "not_helpful"
                )
                _uiState.update {
                    it.copy(feedbackSubmitted = it.feedbackSubmitted + assistantMessageId)
                }
            } catch (_: Exception) {
                _uiState.update {
                    it.copy(error = UiError("Feedback could not be saved.", retryable = true))
                }
            }
        }
    }

    private fun restoreOrCreateSession() {
        requestJob = viewModelScope.launch {
            _uiState.update { it.copy(isSending = true) }
            try {
                val existingSessionId = _uiState.value.sessionId
                if (existingSessionId == null) {
                    ensureSession()
                } else {
                    val session = repository.resumeSession(existingSessionId)
                    _uiState.update {
                        it.copy(
                            sessionId = session.sessionId,
                            messages = session.recentResponses.map { response -> AssistantItem.Response(response) }
                        )
                    }
                }
                _uiState.update { it.copy(isSending = false, error = null) }
            } catch (_: Exception) {
                showSessionUnavailable()
            }
        }
    }

    private suspend fun ensureSession(): String {
        _uiState.value.sessionId?.let { return it }
        val session = repository.createSession(locale = locale, processingConsent = false)
        savedStateHandle[SESSION_ID_KEY] = session.sessionId
        _uiState.update {
            it.copy(
                sessionId = session.sessionId,
                messages = session.recentResponses.map { response -> AssistantItem.Response(response) }
            )
        }
        return session.sessionId
    }

    private fun showSessionUnavailable() {
        _uiState.update {
            it.copy(
                isSending = false,
                error = UiError(
                    message = "Tourtect cannot start an assistant session right now.",
                    retryable = true
                )
            )
        }
    }

    private companion object {
        const val SESSION_ID_KEY = "assistant_session_id"
        const val INPUT_KEY = "assistant_input_draft"
    }
}
