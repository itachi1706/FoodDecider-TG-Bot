import {NextRequest, NextResponse} from "next/server";


export async function middleware(req: NextRequest) {
    // Get auth data from cookie
    console.log("Middleware");
    const authCookie = req.cookies.get("auth");
    if (!authCookie) {
        return NextResponse.redirect("/auth/signin");
    }

    return NextResponse.next();
}