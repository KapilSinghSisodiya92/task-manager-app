"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    // Check if a token exists in local storage
    const token = localStorage.getItem("token");

    if (token) {
      // User is already logged in, send them straight to the workspace
      router.push("/dashboard");
    } else {
      // New or unauthenticated visitor, send them to log in
      router.push("/login");
    }
  }, [router]);

  // Render a clean background loader while the dynamic redirect takes place
  return (
    <div className="flex h-screen w-screen items-center justify-center bg-zinc-50 dark:bg-black">
      <div className="flex flex-col items-center gap-3">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-zinc-300 border-t-zinc-900 dark:border-zinc-700 dark:border-t-zinc-50"></div>
        <p className="text-sm font-medium text-zinc-500 dark:text-zinc-400">Loading your workspace...</p>
      </div>
    </div>
  );
}