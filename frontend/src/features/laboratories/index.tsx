import { Building2, Plus } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'

export function Laboratories() {
  const { t } = useTranslation()

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
            <h2 className='text-2xl font-bold tracking-tight'>
              {t('laboratories.title')}
            </h2>
            <p className='text-muted-foreground'>
              {t('laboratories.description')}
            </p>
          </div>
          <Button>
            <Plus className='mr-2 h-4 w-4' />
            {t('laboratories.addLaboratory')}
          </Button>
        </div>

        {/* Placeholder content - will be replaced with data table */}
        <Card>
          <CardHeader>
            <CardTitle className='flex items-center gap-2'>
              <Building2 className='h-5 w-5' />
              {t('laboratories.laboratoryManagement')}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className='text-muted-foreground'>
              {t('laboratories.comingSoon')}
            </p>
          </CardContent>
        </Card>
      </Main>
    </>
  )
}
