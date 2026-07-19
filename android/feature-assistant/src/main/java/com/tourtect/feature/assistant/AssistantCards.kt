@file:OptIn(androidx.compose.foundation.layout.ExperimentalLayoutApi::class)

package com.tourtect.feature.assistant

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.FlowRow
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.tourtect.core.model.AssistantConfirmation
import com.tourtect.core.model.AssistantEvidence
import com.tourtect.core.model.AssistantIntent
import com.tourtect.core.model.AssistantResponse
import com.tourtect.core.model.AssistantSafetyState
import com.tourtect.core.model.AssistantSuggestedAction

@Composable
internal fun AssistantResponseCard(
    response: AssistantResponse,
    feedbackSubmitted: Boolean,
    onSuggestedAction: (AssistantSuggestedAction) -> Unit,
    onFeedback: (Boolean) -> Unit
) {
    val isSafety = response.safetyState == AssistantSafetyState.CRITICAL ||
        response.safetyState == AssistantSafetyState.URGENT ||
        response.intent == AssistantIntent.SAFETY_ASSESSMENT ||
        response.intent == AssistantIntent.EMERGENCY_HELP
    val isPrice = response.intent == AssistantIntent.PRICE_CHECK ||
        response.intent == AssistantIntent.PRICE_EXPLANATION
    val isTranslation = response.intent == AssistantIntent.TRANSLATION ||
        response.intent == AssistantIntent.LIVE_TRANSLATION

    Card(
        colors = CardDefaults.cardColors(
            containerColor = when {
                isSafety -> MaterialTheme.colorScheme.errorContainer
                isPrice -> MaterialTheme.colorScheme.secondaryContainer
                isTranslation -> MaterialTheme.colorScheme.tertiaryContainer
                else -> MaterialTheme.colorScheme.surfaceVariant
            }
        ),
        modifier = Modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Text(
                text = when {
                    isSafety -> "Safety guidance"
                    isPrice -> "Price insight"
                    isTranslation -> "Translation"
                    response.intent == AssistantIntent.PLACE_DISCOVERY ||
                        response.intent == AssistantIntent.PLACE_INFORMATION -> "Place context"
                    else -> "Tourtect"
                },
                style = MaterialTheme.typography.labelLarge
            )
            Spacer(Modifier.height(6.dp))
            Text(response.message, style = MaterialTheme.typography.bodyLarge)
            Spacer(Modifier.height(8.dp))
            Text(
                "Confidence ${(response.confidence * 100).toInt()}%" +
                    (response.freshness?.let { " · $it" } ?: ""),
                style = MaterialTheme.typography.labelMedium
            )
            if (response.fallbackUsed) {
                Spacer(Modifier.height(8.dp))
                Text(
                    "A deterministic fallback was used.",
                    style = MaterialTheme.typography.labelMedium
                )
            }
            if (response.factorsConsidered.isNotEmpty()) {
                Spacer(Modifier.height(12.dp))
                Text("Why Tourtect is showing this", style = MaterialTheme.typography.titleSmall)
                response.factorsConsidered.forEach { Text("• $it") }
            }
            if (response.missingInformation.isNotEmpty()) {
                ClarificationCard(response.missingInformation)
            }
            if (response.evidence.isNotEmpty()) {
                EvidenceList(response.evidence)
            }
            if (response.suggestedActions.isNotEmpty()) {
                SuggestedActionBar(response.suggestedActions, onSuggestedAction)
            }
            Spacer(Modifier.height(8.dp))
            if (feedbackSubmitted) {
                Text("Feedback saved for review", style = MaterialTheme.typography.labelSmall)
            } else {
                Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    TextButton(onClick = { onFeedback(true) }) { Text("Helpful") }
                    TextButton(onClick = { onFeedback(false) }) { Text("Not helpful") }
                }
            }
        }
    }
}

@Composable
internal fun ClarificationCard(missingInformation: List<String>) {
    Column(modifier = Modifier.padding(top = 12.dp)) {
        Text("More information needed", style = MaterialTheme.typography.titleSmall)
        missingInformation.forEach { Text("• $it", style = MaterialTheme.typography.bodyMedium) }
    }
}

@Composable
internal fun EvidenceList(evidence: List<AssistantEvidence>) {
    Column(modifier = Modifier.padding(top = 12.dp)) {
        Text("Evidence", style = MaterialTheme.typography.titleSmall)
        evidence.forEach { item ->
            Text(item.title, style = MaterialTheme.typography.labelLarge)
            Text(item.summary, style = MaterialTheme.typography.bodyMedium)
            Text(
                "${item.evidenceLevel} · ${item.freshness}",
                style = MaterialTheme.typography.labelSmall
            )
            Spacer(Modifier.height(8.dp))
        }
    }
}

@Composable
internal fun ConfirmationCard(
    confirmation: AssistantConfirmation,
    enabled: Boolean,
    onDecision: (Boolean) -> Unit
) {
    Card(
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.primaryContainer
        ),
        modifier = Modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Text(confirmation.title, style = MaterialTheme.typography.titleMedium)
            Spacer(Modifier.height(4.dp))
            Text(confirmation.description)
            Spacer(Modifier.height(12.dp))
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                Button(enabled = enabled, onClick = { onDecision(true) }) { Text("Confirm") }
                OutlinedButton(enabled = enabled, onClick = { onDecision(false) }) { Text("Not now") }
            }
        }
    }
}

@Composable
internal fun SuggestedActionBar(
    actions: List<AssistantSuggestedAction>,
    onAction: (AssistantSuggestedAction) -> Unit
) {
    FlowRow(
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        modifier = Modifier.padding(top = 12.dp)
    ) {
        actions.forEach { action ->
            OutlinedButton(onClick = { onAction(action) }) { Text(action.label) }
        }
    }
}

@Composable
internal fun ProviderDegradedState(message: String, onRetry: () -> Unit) {
    Card(
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.errorContainer
        ),
        modifier = Modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Text(message, style = MaterialTheme.typography.titleSmall)
            Text("Retry the assistant or open a connected manual fallback if one is available in this release.")
            TextButton(onClick = onRetry) { Text("Retry session") }
        }
    }
}
