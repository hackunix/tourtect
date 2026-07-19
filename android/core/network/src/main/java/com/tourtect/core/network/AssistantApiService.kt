package com.tourtect.core.network

import com.google.gson.annotations.SerializedName
import com.tourtect.core.model.AssistantConfirmation
import com.tourtect.core.model.AssistantEvidence
import com.tourtect.core.model.AssistantIntent
import com.tourtect.core.model.AssistantResponse
import com.tourtect.core.model.AssistantSafetyState
import com.tourtect.core.model.AssistantSession
import com.tourtect.core.model.AssistantSuggestedAction
import com.tourtect.core.model.AssistantToolResult
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.POST
import retrofit2.http.Path

interface AssistantApiService {
    @POST("v1/assistant/sessions")
    suspend fun createSession(
        @Header("X-Request-ID") requestId: String,
        @Body request: CreateAssistantSessionDto
    ): AssistantSessionDto

    @GET("v1/assistant/sessions/{sessionId}")
    suspend fun getSession(
        @Header("X-Request-ID") requestId: String,
        @Path("sessionId") sessionId: String
    ): AssistantSessionDto

    @POST("v1/assistant/sessions/{sessionId}/messages")
    suspend fun sendMessage(
        @Header("X-Request-ID") requestId: String,
        @Path("sessionId") sessionId: String,
        @Body request: AssistantMessageRequestDto
    ): AssistantResponseDto

    @POST("v1/assistant/sessions/{sessionId}/confirmations")
    suspend fun confirmAction(
        @Header("X-Request-ID") requestId: String,
        @Path("sessionId") sessionId: String,
        @Body request: AssistantConfirmationRequestDto
    ): AssistantConfirmationResultDto

    @POST("v1/assistant/sessions/{sessionId}/feedback")
    suspend fun submitFeedback(
        @Header("X-Request-ID") requestId: String,
        @Path("sessionId") sessionId: String,
        @Body request: AssistantFeedbackRequestDto
    ): AssistantFeedbackReceiptDto

    @DELETE("v1/assistant/sessions/{sessionId}")
    suspend fun deleteSession(
        @Header("X-Request-ID") requestId: String,
        @Path("sessionId") sessionId: String
    )
}

data class CreateAssistantSessionDto(
    val locale: String,
    @SerializedName("target_locale") val targetLocale: String? = null,
    @SerializedName("place_id") val placeId: String? = null,
    @SerializedName("approximate_region") val approximateRegion: String? = null,
    @SerializedName("interaction_mode") val interactionMode: String = "text",
    @SerializedName("processing_consent") val processingConsent: Boolean = false
)

data class AssistantMessageRequestDto(
    @SerializedName("message_id") val messageId: String,
    @SerializedName("input_type") val inputType: String,
    val text: String? = null,
    val locale: String? = null,
    @SerializedName("place_id") val placeId: String? = null,
    @SerializedName("user_confirmed") val userConfirmed: Boolean = false,
    @SerializedName("structured_data") val structuredData: Map<String, Any?>? = null
)

data class AssistantSessionDto(
    @SerializedName("session_id") val sessionId: String,
    val version: Int,
    @SerializedName("expires_at") val expiresAt: String,
    @SerializedName("recent_responses") val recentResponses: List<AssistantResponseDto>? = null
)

data class AssistantEvidenceDto(
    @SerializedName("evidence_id") val evidenceId: String,
    @SerializedName("source_type") val sourceType: String,
    @SerializedName("source_id") val sourceId: String,
    val title: String,
    val summary: String,
    @SerializedName("observed_at") val observedAt: String? = null,
    val freshness: String,
    @SerializedName("evidence_level") val evidenceLevel: String,
    @SerializedName("source_url") val sourceUrl: String? = null
)

data class AssistantToolResultDto(
    @SerializedName("tool_result_id") val toolResultId: String,
    @SerializedName("tool_name") val toolName: String,
    val status: String,
    @SerializedName("duration_ms") val durationMs: Long,
    val output: Map<String, Any?> = emptyMap(),
    @SerializedName("error_category") val errorCategory: String? = null
)

data class AssistantConfirmationDto(
    @SerializedName("confirmation_id") val confirmationId: String,
    val action: String,
    val title: String,
    val description: String,
    @SerializedName("expires_at") val expiresAt: String
)

data class AssistantSuggestedActionDto(
    @SerializedName("action_id") val actionId: String,
    val label: String,
    @SerializedName("action_type") val actionType: String,
    val target: String? = null,
    @SerializedName("requires_confirmation") val requiresConfirmation: Boolean
)

data class AssistantResponseDto(
    @SerializedName("assistant_message_id") val assistantMessageId: String,
    val intent: String,
    val message: String,
    val confidence: Double,
    val evidence: List<AssistantEvidenceDto> = emptyList(),
    @SerializedName("tool_results") val toolResults: List<AssistantToolResultDto> = emptyList(),
    @SerializedName("requested_confirmation") val requestedConfirmation: AssistantConfirmationDto? = null,
    @SerializedName("suggested_actions") val suggestedActions: List<AssistantSuggestedActionDto> = emptyList(),
    @SerializedName("safety_state") val safetyState: String,
    @SerializedName("factors_considered") val factorsConsidered: List<String> = emptyList(),
    @SerializedName("missing_information") val missingInformation: List<String> = emptyList(),
    val freshness: String? = null,
    @SerializedName("dataset_version") val datasetVersion: String? = null,
    @SerializedName("fallback_used") val fallbackUsed: Boolean,
    @SerializedName("trace_id") val traceId: String
)

data class AssistantConfirmationRequestDto(
    @SerializedName("confirmation_id") val confirmationId: String,
    val decision: String
)

data class AssistantConfirmationResultDto(
    @SerializedName("confirmation_id") val confirmationId: String,
    val action: String,
    val status: String,
    @SerializedName("executed_at") val executedAt: String,
    @SerializedName("result_id") val resultId: String? = null
)

data class AssistantFeedbackRequestDto(
    @SerializedName("assistant_message_id") val assistantMessageId: String,
    @SerializedName("feedback_type") val feedbackType: String,
    val field: String? = null,
    @SerializedName("original_value") val originalValue: String? = null,
    @SerializedName("corrected_value") val correctedValue: String? = null,
    @SerializedName("consent_to_contribute") val consentToContribute: Boolean = false
)

data class AssistantFeedbackReceiptDto(
    @SerializedName("feedback_id") val feedbackId: String,
    val status: String,
    @SerializedName("created_at") val createdAt: String
)

internal fun AssistantSessionDto.toModel() = AssistantSession(
    sessionId = sessionId,
    version = version,
    expiresAt = expiresAt,
    recentResponses = recentResponses.orEmpty().map { it.toModel() }
)

internal fun AssistantResponseDto.toModel() = AssistantResponse(
    assistantMessageId = assistantMessageId,
    intent = AssistantIntent.fromWire(intent),
    message = message,
    confidence = confidence,
    evidence = evidence.map {
        AssistantEvidence(
            evidenceId = it.evidenceId,
            sourceType = it.sourceType,
            sourceId = it.sourceId,
            title = it.title,
            summary = it.summary,
            observedAt = it.observedAt,
            freshness = it.freshness,
            evidenceLevel = it.evidenceLevel,
            sourceUrl = it.sourceUrl
        )
    },
    toolResults = toolResults.map {
        AssistantToolResult(
            toolResultId = it.toolResultId,
            toolName = it.toolName,
            status = it.status,
            durationMs = it.durationMs,
            output = it.output,
            errorCategory = it.errorCategory
        )
    },
    requestedConfirmation = requestedConfirmation?.let {
        AssistantConfirmation(it.confirmationId, it.action, it.title, it.description, it.expiresAt)
    },
    suggestedActions = suggestedActions.map {
        AssistantSuggestedAction(it.actionId, it.label, it.actionType, it.target, it.requiresConfirmation)
    },
    safetyState = AssistantSafetyState.fromWire(safetyState),
    factorsConsidered = factorsConsidered,
    missingInformation = missingInformation,
    freshness = freshness,
    datasetVersion = datasetVersion,
    fallbackUsed = fallbackUsed,
    traceId = traceId
)
