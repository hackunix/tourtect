import { NextRequest, NextResponse } from "next/server";

export const dynamic = "force-dynamic";

const API_URL = process.env.API_URL || "http://localhost:8080";

async function proxy(request: NextRequest) {
  const url = new URL(request.url);
  const backendUrl = `${API_URL}${url.pathname}${url.search}`;

  const headers = new Headers(request.headers);
  headers.delete("host");

  let body: any = undefined;
  if (request.method !== "GET" && request.method !== "HEAD") {
    try {
      body = await request.arrayBuffer();
    } catch {
      // Ignored if body cannot be parsed or is empty
    }
  }

  try {
    const response = await fetch(backendUrl, {
      method: request.method,
      headers,
      body,
      redirect: "manual",
    });

    const responseHeaders = new Headers(response.headers);
    return new NextResponse(response.body, {
      status: response.status,
      statusText: response.statusText,
      headers: responseHeaders,
    });
  } catch (err: any) {
    console.error(`Failed to proxy request to ${backendUrl}:`, err);
    return NextResponse.json(
      { detail: `Failed to proxy request: ${err.message}` },
      { status: 502 }
    );
  }
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const DELETE = proxy;
export const PATCH = proxy;
export const OPTIONS = proxy;
