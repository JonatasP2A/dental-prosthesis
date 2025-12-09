import { SignedIn, SignedOut } from '@clerk/nextjs';
import Link from 'next/link';

export default function Home() {
  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">
          Welcome to Dental Prosthesis Platform
        </h1>
        <p className="text-lg text-gray-600">
          Manage your laboratory, clients, and orders all in one place.
        </p>
      </div>

      <SignedOut>
        <div className="max-w-md mx-auto bg-white p-8 rounded-lg shadow-md text-center">
          <h2 className="text-2xl font-semibold mb-4">Get Started</h2>
          <p className="text-gray-600 mb-6">
            Sign in or create an account to access the platform.
          </p>
        </div>
      </SignedOut>

      <SignedIn>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Link
            href="/laboratories"
            className="p-6 bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow"
          >
            <h2 className="text-xl font-semibold mb-2">Laboratories</h2>
            <p className="text-gray-600">
              Manage your dental prosthesis laboratories
            </p>
          </Link>

          <Link
            href="/clients"
            className="p-6 bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow"
          >
            <h2 className="text-xl font-semibold mb-2">Clients</h2>
            <p className="text-gray-600">Manage your dental clinic clients</p>
          </Link>

          <Link
            href="/orders"
            className="p-6 bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow"
          >
            <h2 className="text-xl font-semibold mb-2">Orders</h2>
            <p className="text-gray-600">Track and manage prosthesis orders</p>
          </Link>

          <Link
            href="/api-example"
            className="p-6 bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow border-2 border-dashed border-blue-300"
          >
            <h2 className="text-xl font-semibold mb-2">🔑 API Example</h2>
            <p className="text-gray-600">See how to access Clerk JWT token</p>
          </Link>
        </div>
      </SignedIn>
    </div>
  );
}
