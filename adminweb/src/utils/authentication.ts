
export async function authCheck() : Promise<boolean> {
    const authData = localStorage.getItem("authData");

    if (authData) {
        console.log(authData);
        const result = await fetch("/api/auth/validate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: authData,
        });

        if (result.status === 200) {
            return true;
        }
    }

    return false;
}