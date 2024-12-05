import {NextRequest} from "next/server";

export async function GET(req: NextRequest) {
    // TODO: Get user from header to validate
    console.log(req.headers);

    return Response.json({ foodCount: 30, groupCount: 2, locationCount: 3, userCount: 1});
}
