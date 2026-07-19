@file:OptIn(androidx.compose.foundation.layout.ExperimentalLayoutApi::class)

package com.tourtect.feature.assistant

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.FlowRow
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.weight
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.tourtect.core.model.AssistantItem
import com.tourtect.core.model.AssistantSuggestedAction

@Composable
fun AssistantRoute(
    onOpenDestination: (String) -> Unit,
    viewModel: AssistantViewModel = hiltViewModel()
) {
    val state by viewModel.uiState.collectAsState()
    AssistantScreen(
        state = state,
        onInputChange = viewModel::updateInput,
        onSend = viewModel::sendMessage,
        onCancel = viewModel::cancelRequest,
        onReset = viewModel::resetSession,
        onQuickPrompt = viewModel::usePrompt,
        onConfirmation = viewModel::decideConfirmation,
        onFeedback = viewModel::submitFeedback,
        onOpenDestination = onOpenDestination
    )
}

@Composable
fun AssistantScreen(
    state: AssistantUiState,
    onInputChange: (String) -> Unit,
    onSend: () -> Unit,
    onCancel: () -> Unit,
    onReset: () -> Unit,
    onQuickPrompt: (String) -> Unit,
    onConfirmation: (Boolean) -> Unit,
    onFeedback: (String, Boolean) -> Unit,
    onOpenDestination: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    Column(modifier = modifier.fillMaxSize().padding(horizontal = 16.dp)) {
        Row(
            horizontalArrangement = Arrangement.SpaceBetween,
            modifier = Modifier.fillMaxWidth().padding(top = 12.dp)
        ) {
            Column {
                Text("What is happening?", style = MaterialTheme.typography.headlineSmall)
                Text("Speak, type, or show Tourtect what you see.")
            }
            TextButton(onClick = onReset) { Text("Reset") }
        }

        LazyColumn(
            verticalArrangement = Arrangement.spacedBy(10.dp),
            modifier = Modifier.weight(1f).fillMaxWidth().padding(vertical = 12.dp)
        ) {
            if (state.messages.isEmpty()) {
                item { QuickActions(onQuickPrompt, onOpenDestination) }
            }
            state.error?.let { error ->
                item { ProviderDegradedState(error.message, onReset) }
            }
            items(state.messages, key = { it.stableId }) { item ->
                when (item) {
                    is AssistantItem.UserMessage -> UserMessageCard(item.text)
                    is AssistantItem.Response -> AssistantResponseCard(
                        response = item.value,
                        feedbackSubmitted = item.stableId in state.feedbackSubmitted,
                        onSuggestedAction = { action -> routeAction(action, onOpenDestination, onQuickPrompt) },
                        onFeedback = { onFeedback(item.stableId, it) }
                    )
                    is AssistantItem.StatusNotice -> StatusNoticeCard(item.title, item.description)
                }
            }
            state.pendingConfirmation?.let { confirmation ->
                item {
                    ConfirmationCard(
                        confirmation = confirmation,
                        enabled = !state.isSending,
                        onDecision = onConfirmation
                    )
                }
            }
        }

        OutlinedTextField(
            value = state.input,
            onValueChange = onInputChange,
            enabled = !state.isSending,
            label = { Text("Type a message") },
            modifier = Modifier.fillMaxWidth()
        )
        Row(
            horizontalArrangement = Arrangement.spacedBy(8.dp),
            modifier = Modifier.fillMaxWidth().padding(vertical = 10.dp)
        ) {
            OutlinedButton(onClick = { onOpenDestination("live") }) { Text("Hold to speak") }
            OutlinedButton(onClick = { onOpenDestination("lens") }) { Text("Show Tourtect") }
            if (state.isSending) {
                OutlinedButton(onClick = onCancel) { Text("Cancel") }
                CircularProgressIndicator(modifier = Modifier.height(32.dp))
            } else {
                Button(onClick = onSend, enabled = state.input.isNotBlank()) { Text("Send") }
            }
        }
    }
}

@Composable
private fun QuickActions(
    onQuickPrompt: (String) -> Unit,
    onOpenDestination: (String) -> Unit
) {
    Column {
        Text("Quick actions", style = MaterialTheme.typography.titleMedium)
        FlowRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            OutlinedButton(onClick = { onQuickPrompt("Check whether this price is reasonable: ") }) {
                Text("Check a price")
            }
            OutlinedButton(onClick = { onQuickPrompt("Translate this conversation: ") }) {
                Text("Translate")
            }
            OutlinedButton(onClick = { onOpenDestination("lens") }) { Text("Scan menu or receipt") }
            OutlinedButton(onClick = { onQuickPrompt("I am concerned about this situation: ") }) {
                Text("Describe a suspicious situation")
            }
            OutlinedButton(onClick = { onQuickPrompt("What should I know about this place: ") }) {
                Text("Ask about a place")
            }
            OutlinedButton(onClick = { onOpenDestination("sos") }) { Text("Emergency help") }
        }
    }
}

@Composable
private fun UserMessageCard(text: String) {
    Row(horizontalArrangement = Arrangement.End, modifier = Modifier.fillMaxWidth()) {
        Card { Text(text, modifier = Modifier.padding(14.dp)) }
    }
}

@Composable
private fun StatusNoticeCard(title: String, description: String) {
    Card(modifier = Modifier.fillMaxWidth()) {
        Column(modifier = Modifier.padding(14.dp)) {
            Text(title, style = MaterialTheme.typography.titleSmall)
            Spacer(Modifier.height(4.dp))
            Text(description)
        }
    }
}

private fun routeAction(
    action: AssistantSuggestedAction,
    onOpenDestination: (String) -> Unit,
    onQuickPrompt: (String) -> Unit
) {
    when (action.actionType) {
        "manual_price_check" -> onOpenDestination("price-check")
        "manual_safety_assessment" -> onOpenDestination("safety")
        "offline_directory" -> onOpenDestination("sos")
        "deep_link" -> action.target?.let(onOpenDestination)
        "clarify" -> onQuickPrompt(action.label)
    }
}
