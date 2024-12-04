"use client";
import React from "react";
import {LoginButton} from "@telegram-auth/react";

export default function Signin() {
    const botUsername = process.env.NEXT_PUBLIC_BOT_USERNAME || "";
    console.log("Bot Username:", botUsername);
    return (
        <>
            <div className="my-6 flex items-center justify-center">
                <LoginButton
                    botUsername={botUsername}
                    buttonSize={"large"}
                    showAvatar={true}
                    onAuthCallback={(authData) => {
                        console.log(authData);

                        // Store the auth data in the local storage
                        localStorage.setItem("authData", JSON.stringify(authData));

                        // Redirect to the dashboard
                        window.location.href = "/";
                    }}
                />
            </div>
        </>
    );
}
