import {NextRequest, NextResponse} from "next/server";

export function middleware(req: NextRequest) {
    // Get URL path from URL
    const uri = new URL(req.url);
    // Exclude auth routes and sign in page
    if (uri.pathname === "/auth/signin" || uri.pathname.startsWith("/api/auth") || uri.pathname == "/") {
        return NextResponse.next();
    }
    console.log("Checking authentication on path", uri.pathname)

    const authCookie = req.cookies.get("auth");
    if (!authCookie) {
        console.log("Redirecting to sign in on path", uri.pathname)

        const nextUri = req.nextUrl.clone()
        nextUri.pathname = '/auth/signin'
        return NextResponse.redirect(nextUri);
    }

    return NextResponse.next();
}

export const config = {
    matcher: ["/api/(.*)", "/((?!api|_next/static|_next/image|images|favicon.ico|sitemap.xml|robots.txt).*)"],
}