import { HardHat, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'

export function Technicians() {
  return (
    <>
      <Header fixed>
        <Search />
        <div className='ms-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      <Main className='flex flex-1 flex-col gap-4 sm:gap-6'>
        <div className='flex flex-wrap items-end justify-between gap-2'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>Technicians</h2>
            <p className='text-muted-foreground'>
              Manage your laboratory technicians.
            </p>
          </div>
          <Button>
            <Plus className='mr-2 h-4 w-4' />
            Add Technician
          </Button>
        </div>

        {/* Placeholder content - will be replaced with data table */}
        <Card>
          <CardHeader>
            <CardTitle className='flex items-center gap-2'>
              <HardHat className='h-5 w-5' />
              Technician Management
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className='text-muted-foreground'>
              Technician management interface coming soon. This page will allow
              you to manage lab staff, their roles, and specializations.
            </p>
          </CardContent>
        </Card>
      </Main>
    </>
  )
}
