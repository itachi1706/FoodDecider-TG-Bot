"use client";
import React from "react";
import {LoginButton} from "@telegram-auth/react";
import ErrorAlert from "@/components/Alerts/ErrorAlert";

export default function Signin() {

  const [showError, setShowError] = React.useState(false);

    const botUsername = process.env.NEXT_PUBLIC_BOT_USERNAME || "";
    console.log("Bot Username:", botUsername);
    return (
        <>
          {showError &&
            <ErrorAlert
              header="Unauthorized Telegram User"
              message="You must be a recognized administrator to use this website" />
          }
            <div className="my-6 flex items-center justify-center">
                <LoginButton
                    botUsername={botUsername}
                    buttonSize={"large"}
                    showAvatar={true}
                    onAuthCallback={async (authData) => {
                        console.log(authData);

                      // Store the auth data in the local storage
                      localStorage.setItem("authData", JSON.stringify(authData));
                        // Login user
                        try {
                            const resp = await fetch("/api/auth/login", {
                                method: "POST",
                                headers: {
                                    "Content-Type": "application/json",
                                },
                                body: JSON.stringify(authData),
                            });
                            if (resp.status === 403) {
                              // Not admin
                              setShowError(true);
                              return;

                            } else if (resp.status !== 200) {
                                console.error("Failed to login");
                                return;
                            }
                            // Redirect to the dashboard
                            window.location.href = "/";
                        } catch (error) {
                            console.error("Failed to login", error);
                        }
                    }}
                />
            </div>
        </>
    );
}
