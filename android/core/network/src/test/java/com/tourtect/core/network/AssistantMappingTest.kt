package com.tourtect.core.network

import com.tourtect.core.model.AssistantIntent
import com.tourtect.core.model.AssistantSafetyState
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Test

class AssistantMappingTest {
    @Test
    fun `response mapping preserves deterministic and provenance fields`() {
        val response = AssistantResponseDto(
            assistantMessageId = "message-1",
            intent = "price_check",
            message = "Structured result",
            confidence = 0.88,
            evidence = listOf(
                AssistantEvidenceDto(
                    evidenceId = "evidence-1",
                    sourceType = "price_snapshot",
                    sourceId = "snapshot-1",
                    title = "Reference range",
                    summary = "Verified observations",
                    freshness = "fresh",
                    evidenceLevel = "verified"
                )
            ),
            toolResults = listOf(
                AssistantToolResultDto(
                    toolResultId = "tool-1",
                    toolName = "evaluate_price",
                    status = "succeeded",
                    durationMs = 12,
                    output = mapOf("alert_level" to "high")
                )
            ),
            safetyState = "non_emergency",
            factorsConsidered = listOf("Entered price"),
            datasetVersion = "price-v1",
            fallbackUsed = false,
            traceId = "trace-1"
        ).toModel()

        assertEquals(AssistantIntent.PRICE_CHECK, response.intent)
        assertEquals(AssistantSafetyState.NON_EMERGENCY, response.safetyState)
        assertEquals("evidence-1", response.evidence.single().evidenceId)
        assertEquals("high", response.toolResults.single().output["alert_level"])
        assertEquals("price-v1", response.datasetVersion)
        assertFalse(response.fallbackUsed)
    }
}
