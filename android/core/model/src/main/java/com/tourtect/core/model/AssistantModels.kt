package com.tourtect.core.model

enum class AssistantIntent {
    GENERAL_TRAVEL_QUESTION,
    PLACE_DISCOVERY,
    PLACE_INFORMATION,
    PRICE_CHECK,
    PRICE_EXPLANATION,
    TRANSLATION,
    LIVE_TRANSLATION,
    MENU_OR_RECEIPT_ANALYSIS,
    SCAM_PATTERN_ASSESSMENT,
    SAFETY_ASSESSMENT,
    EMERGENCY_HELP,
    COMMUNITY_SEARCH,
    CREATE_REPORT_DRAFT,
    UNKNOWN;

    companion object {
        fun fromWire(value: String): AssistantIntent = entries.firstOrNull {
            it.name.equals(value, ignoreCase = true)
        } ?: UNKNOWN
    }
}

enum class AssistantSafetyState {
    CRITICAL,
    URGENT,
    NON_EMERGENCY,
    INFORMATION,
    UNKNOWN;

    companion object {
        fun fromWire(value: String): AssistantSafetyState = entries.firstOrNull {
            it.name.equals(value, ignoreCase = true)
        } ?: UNKNOWN
    }
}

data class AssistantEvidence(
    val evidenceId: String,
    val sourceType: String,
    val sourceId: String,
    val title: String,
    val summary: String,
    val observedAt: String?,
    val freshness: String,
    val evidenceLevel: String,
    val sourceUrl: String?
)

data class AssistantToolResult(
    val toolResultId: String,
    val toolName: String,
    val status: String,
    val durationMs: Long,
    val output: Map<String, Any?>,
    val errorCategory: String?
)

data class AssistantConfirmation(
    val confirmationId: String,
    val action: String,
    val title: String,
    val description: String,
    val expiresAt: String
)

data class AssistantSuggestedAction(
    val actionId: String,
    val label: String,
    val actionType: String,
    val target: String?,
    val requiresConfirmation: Boolean
)

data class AssistantResponse(
    val assistantMessageId: String,
    val intent: AssistantIntent,
    val message: String,
    val confidence: Double,
    val evidence: List<AssistantEvidence>,
    val toolResults: List<AssistantToolResult>,
    val requestedConfirmation: AssistantConfirmation?,
    val suggestedActions: List<AssistantSuggestedAction>,
    val safetyState: AssistantSafetyState,
    val factorsConsidered: List<String>,
    val missingInformation: List<String>,
    val freshness: String?,
    val datasetVersion: String?,
    val fallbackUsed: Boolean,
    val traceId: String
)

data class AssistantSession(
    val sessionId: String,
    val version: Int,
    val expiresAt: String,
    val recentResponses: List<AssistantResponse>
)

sealed interface AssistantItem {
    val stableId: String

    data class UserMessage(
        override val stableId: String,
        val text: String
    ) : AssistantItem

    data class Response(
        val value: AssistantResponse
    ) : AssistantItem {
        override val stableId: String = value.assistantMessageId
    }

    data class StatusNotice(
        override val stableId: String,
        val title: String,
        val description: String
    ) : AssistantItem
}
