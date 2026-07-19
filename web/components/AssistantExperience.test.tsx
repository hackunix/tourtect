import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import AssistantExperience from "./AssistantExperience";
import { api, AssistantResponse, ProblemDetailError } from "@/lib/api";

vi.mock("@/lib/api",async(importOriginal)=>{
  const actual=await importOriginal<typeof import("@/lib/api")>();
  return{...actual,api:{...actual.api,createAssistantSession:vi.fn(),getAssistantSession:vi.fn(),deleteAssistantSession:vi.fn(),createAssistantMessage:vi.fn(),confirmAssistantAction:vi.fn(),createAssistantFeedback:vi.fn()}}
});

const session={session_id:"11111111-1111-4111-8111-111111111111",version:1,created_at:"2026-07-19T00:00:00Z",updated_at:"2026-07-19T00:00:00Z",expires_at:"2026-07-19T01:00:00Z",context:{locale:"en",interaction_mode:"text",consent_state:{processing:true,contribution:false,publish:false}}};
const response:AssistantResponse={
  assistant_message_id:"22222222-2222-4222-8222-222222222222",intent:"price_check",message:"The entered fare is above the verified reference range.",confidence:.88,
  evidence:[{evidence_id:"33333333-3333-4333-8333-333333333333",source_type:"price_snapshot",source_id:"snapshot-1",title:"Airport taxi range",summary:"Based on 37 recent observations.",freshness:"fresh",evidence_level:"verified"}],
  tool_results:[{tool_result_id:"44444444-4444-4444-8444-444444444444",tool_name:"evaluate_price",status:"succeeded",duration_ms:12,output:{note:"deterministic result preserved"}}],
  requested_confirmation:{confirmation_id:"55555555-5555-4555-8555-555555555555",action:"create_report_draft",title:"Create a private report draft?",description:"Only confirmed facts will be included.",expires_at:"2026-07-19T01:00:00Z"},
  suggested_actions:[],safety_state:"non_emergency",factors_considered:["Entered price: 900,000 VND"],missing_information:[],fallback_used:false,trace_id:"trace-1",
};

describe("AssistantExperience",()=>{
  beforeEach(()=>{vi.clearAllMocks();localStorage.clear();vi.mocked(api.createAssistantSession).mockResolvedValue(session);vi.mocked(api.deleteAssistantSession).mockResolvedValue(undefined);vi.mocked(api.confirmAssistantAction).mockResolvedValue({confirmation_id:response.requested_confirmation!.confirmation_id,action:"create_report_draft",status:"confirmed",executed_at:"2026-07-19T00:05:00Z"});vi.mocked(api.createAssistantFeedback).mockResolvedValue({feedback_id:"feedback-1",status:"quarantined",created_at:"2026-07-19T00:05:00Z"})});
  afterEach(cleanup);

  it("creates a consented session and renders evidence plus server-issued confirmation",async()=>{
    vi.mocked(api.createAssistantMessage).mockResolvedValue(response);
    render(<AssistantExperience locale="en"/>);
    await screen.findByText("Ready when you are.");
    fireEvent.change(screen.getByLabelText("Describe what is happening"),{target:{value:"The driver wants 900,000 VND."}});
    fireEvent.click(screen.getByRole("button",{name:"Send"}));
    expect(await screen.findByText(response.message)).toBeInTheDocument();
    expect(screen.getByRole("region",{name:"Evidence used"})).toHaveTextContent("37 recent observations");
    expect(api.createAssistantSession).toHaveBeenCalledWith(expect.objectContaining({processing_consent:true}),expect.any(AbortSignal));
    expect(api.createAssistantMessage).toHaveBeenCalledWith(session.session_id,expect.objectContaining({input_type:"text",text:"The driver wants 900,000 VND."}),expect.any(AbortSignal));
    fireEvent.click(screen.getByRole("button",{name:"Confirm"}));
    await waitFor(()=>expect(api.confirmAssistantAction).toHaveBeenCalledWith(session.session_id,{confirmation_id:response.requested_confirmation!.confirmation_id,decision:"confirmed"}));
    expect(await screen.findByText(/Action confirmed by the backend/)).toBeInTheDocument();
  });

  it("shows deterministic fallback options instead of a fake assistant answer",async()=>{
    vi.mocked(api.createAssistantMessage).mockRejectedValue(new ProblemDetailError(503,"Provider unavailable","request-1"));
    render(<AssistantExperience locale="en"/>);
    await screen.findByText("Ready when you are.");
    fireEvent.change(screen.getByLabelText("Describe what is happening"),{target:{value:"Please translate this."}});
    fireEvent.click(screen.getByRole("button",{name:"Send"}));
    expect(await screen.findByRole("heading",{name:"Automatic understanding is unavailable"})).toBeInTheDocument();
    expect(screen.getByText(/Provider unavailable/)).toBeInTheDocument();
    expect(screen.queryByText(response.message)).not.toBeInTheDocument();
  });
});
