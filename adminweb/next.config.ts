import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    /* config options here */
    images: {
        remotePatterns: [
            {
                protocol: "http",
                hostname: "localhost",
            },
            {
                protocol: "https",
                hostname: "t.me",
            },
            {
                protocol: "https",
                hostname: "cdn.sanity.io",
                port: ""
            },
            {
                protocol: "https",
                hostname: "lh3.googleusercontent.com",
                port: ""
            },
            {
                protocol: "https",
                hostname: "avatars.githubusercontent.com",
                port: ""
            },
            {
                protocol: "https",
                hostname: "pub-b7fd9c30cdbf439183b75041f5f71b92.r2.dev",
                port: ""
            }
        ]
}
};

export default nextConfig;
