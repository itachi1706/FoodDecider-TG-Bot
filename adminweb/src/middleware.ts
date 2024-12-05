import {NextRequest, NextResponse} from "next/server";

export function middleware(req: NextRequest) {
    // Get URL path from URL
    const uri = new URL(req.url);
    const failUri = req.nextUrl.clone()
    failUri.pathname = '/auth/signin'

    // Exclude auth routes and sign in page
    if (uri.pathname === "/auth/signin" || uri.pathname.startsWith("/api/auth") || uri.pathname == "/") {
        return NextResponse.next();
    }
    console.log("Checking authentication on path", uri.pathname)

    const authCookie = req.cookies.get("auth");
    if (!authCookie) {
        console.log("Redirecting to sign in on path", uri.pathname)
        return NextResponse.redirect(failUri);
    }

    // Convert back to JSON
    try {
        const authData = JSON.parse(Buffer.from(authCookie.value, 'base64').toString('utf-8'));
        console.log("Auth data", authData);

        // TODO: Make sure is authorized user (must be admin), else throw to sign in page with error


    } catch (err) {
        console.error("Error parsing auth cookie, rejecting auth", err);
        req.cookies.delete("auth");
        return NextResponse.redirect(failUri);
    }

    return NextResponse.next();
}

export const config = {
    matcher: ["/api/(.*)", "/((?!api|_next/static|_next/image|images|favicon.ico|sitemap.xml|robots.txt).*)"],
}