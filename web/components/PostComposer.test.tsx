import { fireEvent,render,screen,waitFor } from "@testing-library/react";
import { beforeEach,describe,expect,it,vi } from "vitest";
import PostComposer from "./PostComposer";
import { api } from "@/lib/api";

vi.mock("next/navigation",()=>({useRouter:()=>({refresh:vi.fn()})}));
vi.mock("@/lib/api",async(importOriginal)=>{const actual=await importOriginal<typeof import("@/lib/api")>();return{...actual,api:{...actual.api,createDraft:vi.fn(),publishPost:vi.fn()}}});

const draft={post_id:"draft-1",author_id:"author-1",post_type:"discussion",original_locale:"vi-VN",title:"Một kinh nghiệm",body:"Nội dung hữu ích",evidence_level:"none",moderation_status:"draft",created_at:new Date(0).toISOString()};
describe("PostComposer",()=>{beforeEach(()=>vi.clearAllMocks());it("saves a draft and requires a separate publish confirmation",async()=>{vi.mocked(api.createDraft).mockResolvedValue(draft);vi.mocked(api.publishPost).mockResolvedValue({...draft,moderation_status:"published"});render(<PostComposer locale="vi-VN" places={[]} prompt="Bạn muốn du khách biết điều gì?"/>);fireEvent.click(screen.getByRole("button",{name:/Bạn muốn/}));fireEvent.change(screen.getByLabelText("Tiêu đề"),{target:{value:draft.title}});fireEvent.change(screen.getByLabelText("Nội dung"),{target:{value:draft.body}});fireEvent.click(screen.getByRole("button",{name:"Lưu bản nháp"}));await screen.findByText("Xác nhận bản nháp");expect(api.createDraft).toHaveBeenCalledTimes(1);expect(api.publishPost).not.toHaveBeenCalled();fireEvent.click(screen.getByRole("button",{name:"Xuất bản"}));await waitFor(()=>expect(api.publishPost).toHaveBeenCalledWith("draft-1"))})});
