
export async function authCheck() : Promise<boolean> {
    const result = await fetch("/api/auth/validate", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({}),
    });

    return result.status === 200;
}