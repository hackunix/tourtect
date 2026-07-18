import { render,screen } from "@testing-library/react";
import { describe,expect,it } from "vitest";
import FeedHeader from "./FeedHeader";

const dict={following:"Đang theo dõi",nearby:"Gần đây",latest:"Mới nhất",trending:"Nổi bật",safety:"An toàn"} as never;
describe("FeedHeader",()=>{it("marks the selected tab without relying on color",()=>{render(<FeedHeader mode="safety" dict={dict}/>);expect(screen.getByRole("link",{name:/An toàn/})).toHaveAttribute("aria-current","page");expect(screen.getAllByRole("link")).toHaveLength(5)})});
