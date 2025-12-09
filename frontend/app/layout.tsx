import type { Metadata } from 'next';
import {
  ClerkProvider,
  SignInButton,
  SignUpButton,
  SignedIn,
  SignedOut,
  UserButton,
} from '@clerk/nextjs';
import './globals.css';

export const metadata: Metadata = {
  title: 'Dental Prosthesis Platform',
  description: 'Manage dental prosthesis laboratories, clients, and orders',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ClerkProvider>
      <html lang="en">
        <body>
          <header className="border-b">
            <div className="container mx-auto px-4 py-4 flex justify-between items-center">
              <h1 className="text-xl font-bold">Dental Prosthesis Platform</h1>
              <nav>
                <SignedOut>
                  <div className="flex gap-2">
                    <SignInButton mode="modal">
                      <button className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">
                        Sign In
                      </button>
                    </SignInButton>
                    <SignUpButton mode="modal">
                      <button className="px-4 py-2 text-sm font-medium bg-blue-600 text-white rounded hover:bg-blue-700">
                        Sign Up
                      </button>
                    </SignUpButton>
                  </div>
                </SignedOut>
                <SignedIn>
                  <UserButton />
                </SignedIn>
              </nav>
            </div>
          </header>
          <main className="container mx-auto px-4 py-8">{children}</main>
        </body>
      </html>
    </ClerkProvider>
  );
}
