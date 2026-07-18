import { render,screen } from "@testing-library/react";
import { describe,expect,it } from "vitest";
import ProblemDetailsState from "./ProblemDetailsState";
describe("ProblemDetailsState",()=>{it("shows a request id and retry action",()=>{render(<ProblemDetailsState detail="Backend unavailable" requestId="req-123"/>);expect(screen.getByRole("alert")).toHaveTextContent("Backend unavailable");expect(screen.getByText(/req-123/)).toBeInTheDocument();expect(screen.getByRole("link",{name:"Thử lại"})).toHaveAttribute("href","/")})});
