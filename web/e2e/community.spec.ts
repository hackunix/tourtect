import { test,expect } from "@playwright/test";
import AxeBuilder from "@axe-core/playwright";

test("home opens the assistant as the primary experience",async({page})=>{await page.goto("/");await expect(page).toHaveURL(/\/assistant$/);await expect(page.getByRole("heading",{name:"What is happening?"})).toBeVisible();await expect(page.getByRole("navigation",{name:"Mobile navigation"}).or(page.getByRole("navigation",{name:"Primary navigation"}))).toBeVisible()});
test("community shell remains usable with live data or a designed API error",async({page})=>{await page.goto("/community");await expect(page.getByRole("heading",{name:/Cộng đồng du lịch|Traveler community/})).toBeVisible();await expect(page.locator("main")).toBeVisible()});
test("search surface has no automatic accessibility violations",async({page})=>{await page.goto("/search");const results=await new AxeBuilder({page}).analyze();expect(results.violations).toEqual([])});

test("sending a message in the assistant experience", async ({ page }) => {
  await page.goto("/assistant");
  await expect(page.getByText("Ready when you are.")).toBeVisible();
  
  await page.getByPlaceholder("Describe what is happening…").fill("Hello assistant, this is a test");
  
  const sendBtn = page.getByRole("button", { name: "Send" });
  await expect(sendBtn).toBeEnabled();
  await sendBtn.click();
  
  // Verify that the user's message is added to the UI chat stream
  await expect(page.getByText("Hello assistant, this is a test")).toBeVisible();
  
  // Verify that the assistant responds and renders the message article
  await expect(page.locator("article").filter({ hasText: "Tourtect" }).first()).toBeVisible({ timeout: 10000 });
});


