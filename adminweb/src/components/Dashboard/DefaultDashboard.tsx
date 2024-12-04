"use client";
import React, {useEffect, useState} from "react";
import DataStatsBot from "@/components/DataStats/DataStatsBot";
import {authCheck} from "@/utils/authentication";

const DefaultDashboard: React.FC = () => {

    const [authenticated, setAuthenticated] = useState(false);

    // Check if user is authenticated
    useEffect(() => {
        authCheck().then(logged => setAuthenticated(logged)).catch(console.error);
    });

    if (!authenticated) {
        return (<></>)
    }

    return (
        <>
            <DataStatsBot/>
            <div className="mt-4 grid grid-cols-12 gap-4 md:mt-6 md:gap-6 2xl:mt-9 2xl:gap-7.5">
            </div>
        </>
    );
};

export default DefaultDashboard;
