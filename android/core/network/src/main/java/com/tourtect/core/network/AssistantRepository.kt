package com.tourtect.core.network

import com.tourtect.core.model.AssistantResponse
import com.tourtect.core.model.AssistantSession
import java.util.UUID
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AssistantRepository @Inject constructor(
    private val api: AssistantApiService
) {
    suspend fun createSession(
        locale: String,
        targetLocale: String? = null,
        processingConsent: Boolean = false
    ): AssistantSession = api.createSession(
        requestId = UUID.randomUUID().toString(),
        request = CreateAssistantSessionDto(
            locale = locale,
            targetLocale = targetLocale,
            processingConsent = processingConsent
        )
    ).toModel()

    suspend fun resumeSession(sessionId: String): AssistantSession = api.getSession(
        requestId = UUID.randomUUID().toString(),
        sessionId = sessionId
    ).toModel()

    suspend fun sendText(
        sessionId: String,
        text: String,
        locale: String,
        messageId: String = UUID.randomUUID().toString()
    ): AssistantResponse = api.sendMessage(
        requestId = UUID.randomUUID().toString(),
        sessionId = sessionId,
        request = AssistantMessageRequestDto(
            messageId = messageId,
            inputType = "text",
            text = text,
            locale = locale
        )
    ).toModel()

    suspend fun confirmAction(
        sessionId: String,
        confirmationId: String,
        confirmed: Boolean
    ): AssistantConfirmationResultDto = api.confirmAction(
        requestId = UUID.randomUUID().toString(),
        sessionId = sessionId,
        request = AssistantConfirmationRequestDto(
            confirmationId = confirmationId,
            decision = if (confirmed) "confirmed" else "rejected"
        )
    )

    suspend fun submitFeedback(
        sessionId: String,
        assistantMessageId: String,
        feedbackType: String
    ): AssistantFeedbackReceiptDto = api.submitFeedback(
        requestId = UUID.randomUUID().toString(),
        sessionId = sessionId,
        request = AssistantFeedbackRequestDto(
            assistantMessageId = assistantMessageId,
            feedbackType = feedbackType
        )
    )

    suspend fun deleteSession(sessionId: String) {
        api.deleteSession(UUID.randomUUID().toString(), sessionId)
    }
}
