import { test,expect } from "@playwright/test";
import AxeBuilder from "@axe-core/playwright";

test("community shell remains usable with live data or a designed API error",async({page})=>{await page.goto("/");await expect(page.getByRole("heading",{name:/Cộng đồng du lịch|Traveler community/})).toBeVisible();await expect(page.locator("main")).toBeVisible();await expect(page.getByRole("navigation",{name:"Mobile navigation"}).or(page.getByRole("navigation",{name:"Primary navigation"}))).toBeVisible()});
test("search surface has no automatic accessibility violations",async({page})=>{await page.goto("/search");const results=await new AxeBuilder({page}).analyze();expect(results.violations).toEqual([])});
