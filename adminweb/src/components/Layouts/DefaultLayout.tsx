"use client";
import React, {useState, useEffect, useMemo} from "react";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import {authCheck} from "@/utils/authentication";

export default function DefaultLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const [authenticated, setAuthenticated] = useState(false);

  // Set refreshAuth
  const refreshAuth = async () => {
    authCheck().then(logged => setAuthenticated(logged)).catch(console.error);
  }

  // Check if user is authenticated
  useMemo(() => {
    // console.log("Checking auth 3");
    authCheck().then(logged => setAuthenticated(logged)).catch(console.error);
  }, []);


  return (
    <>
      {/* <!-- ===== Page Wrapper Star ===== --> */}
      <div className="flex h-screen overflow-hidden">
        {/* <!-- ===== Sidebar Star ===== --> */}
        <Sidebar sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} isAuthenticated={authenticated} />
        {/* <!-- ===== Sidebar End ===== --> */}

        {/* <!-- ===== Content Area Star ===== --> */}
        <div className="relative flex flex-1 flex-col overflow-y-auto overflow-x-hidden">
          {/* <!-- ===== Header Star ===== --> */}
          <Header sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} isAuthenticated={authenticated} refreshAuth={refreshAuth} />
          {/* <!-- ===== Header End ===== --> */}

          {/* <!-- ===== Main Content Star ===== --> */}
          <main>
            <div className="mx-auto max-w-screen-2xl p-4 md:p-6 2xl:p-10">
              {children}
            </div>
          </main>
          {/* <!-- ===== Main Content End ===== --> */}
        </div>
        {/* <!-- ===== Content Area End ===== --> */}
      </div>
      {/* <!-- ===== Page Wrapper End ===== --> */}
    </>
  );
}
