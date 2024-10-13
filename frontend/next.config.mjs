/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  images: {
    formats: ["image/webp"],
  },
  poweredByHeader: false,
};

export default nextConfig;
