import { defineConfig,devices } from "@playwright/test";
const port=process.env.PLAYWRIGHT_PORT||"3000";
const baseURL=`http://127.0.0.1:${port}`;
export default defineConfig({testDir:"./e2e",fullyParallel:true,use:{baseURL,trace:"retain-on-failure"},projects:[{name:"desktop",use:{...devices["Desktop Chrome"]}},{name:"mobile",use:{...devices["Pixel 7"]}}],webServer:{command:`npm run dev -- --hostname 127.0.0.1 --port ${port}`,url:baseURL,reuseExistingServer:!process.env.CI&&!process.env.PLAYWRIGHT_PORT,timeout:120000}});
