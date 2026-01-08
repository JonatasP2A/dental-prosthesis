import { Building2, ClipboardList, Clock, Package, Users } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Overview } from './components/overview'

export function Dashboard() {
  const { t } = useTranslation()

  return (
    <>
      {/* ===== Top Heading ===== */}
      <Header>
        <Search />
        <div className='ms-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      {/* ===== Main ===== */}
      <Main>
        <div className='mb-4 flex items-center justify-between'>
          <div>
            <h1 className='text-2xl font-bold tracking-tight'>
              {t('dashboard.title')}
            </h1>
            <p className='text-muted-foreground'>{t('dashboard.welcome')}</p>
          </div>
        </div>

        <div className='space-y-4'>
          {/* Stats Cards */}
          <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
            <Card>
              <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                <CardTitle className='text-sm font-medium'>
                  {t('dashboard.activeOrders')}
                </CardTitle>
                <ClipboardList className='text-muted-foreground h-4 w-4' />
              </CardHeader>
              <CardContent>
                <div className='text-2xl font-bold'>24</div>
                <p className='text-muted-foreground text-xs'>
                  8 {t('dashboard.inProduction')}, 6 {t('dashboard.ready')}
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                <CardTitle className='text-sm font-medium'>
                  {t('dashboard.totalClients')}
                </CardTitle>
                <Users className='text-muted-foreground h-4 w-4' />
              </CardHeader>
              <CardContent>
                <div className='text-2xl font-bold'>45</div>
                <p className='text-muted-foreground text-xs'>
                  +3 {t('dashboard.newThisMonth')}
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                <CardTitle className='text-sm font-medium'>
                  {t('dashboard.prosthesisTypes')}
                </CardTitle>
                <Package className='text-muted-foreground h-4 w-4' />
              </CardHeader>
              <CardContent>
                <div className='text-2xl font-bold'>12</div>
                <p className='text-muted-foreground text-xs'>
                  {t('dashboard.crownsAndMore')}
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                <CardTitle className='text-sm font-medium'>
                  {t('dashboard.avgTurnaround')}
                </CardTitle>
                <Clock className='text-muted-foreground h-4 w-4' />
              </CardHeader>
              <CardContent>
                <div className='text-2xl font-bold'>
                  4.2 {t('dashboard.days')}
                </div>
                <p className='text-muted-foreground text-xs'>
                  -0.5 {t('dashboard.daysFromLastMonth')}
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Charts and Recent Activity */}
          <div className='grid grid-cols-1 gap-4 lg:grid-cols-7'>
            <Card className='col-span-1 lg:col-span-4'>
              <CardHeader>
                <CardTitle>{t('dashboard.ordersOverview')}</CardTitle>
                <CardDescription>
                  {t('dashboard.monthlyOrderVolume')}
                </CardDescription>
              </CardHeader>
              <CardContent className='ps-2'>
                <Overview />
              </CardContent>
            </Card>
            <Card className='col-span-1 lg:col-span-3'>
              <CardHeader>
                <CardTitle>{t('dashboard.orderStatus')}</CardTitle>
                <CardDescription>
                  {t('dashboard.currentOrderDistribution')}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className='space-y-4'>
                  <div className='flex items-center justify-between'>
                    <div className='flex items-center gap-2'>
                      <div className='h-3 w-3 rounded-full bg-yellow-500' />
                      <span className='text-sm'>{t('dashboard.received')}</span>
                    </div>
                    <span className='text-sm font-medium'>10</span>
                  </div>
                  <div className='flex items-center justify-between'>
                    <div className='flex items-center gap-2'>
                      <div className='h-3 w-3 rounded-full bg-blue-500' />
                      <span className='text-sm'>
                        {t('dashboard.inProductionStatus')}
                      </span>
                    </div>
                    <span className='text-sm font-medium'>8</span>
                  </div>
                  <div className='flex items-center justify-between'>
                    <div className='flex items-center gap-2'>
                      <div className='h-3 w-3 rounded-full bg-purple-500' />
                      <span className='text-sm'>
                        {t('dashboard.qualityCheck')}
                      </span>
                    </div>
                    <span className='text-sm font-medium'>2</span>
                  </div>
                  <div className='flex items-center justify-between'>
                    <div className='flex items-center gap-2'>
                      <div className='h-3 w-3 rounded-full bg-green-500' />
                      <span className='text-sm'>
                        {t('dashboard.readyStatus')}
                      </span>
                    </div>
                    <span className='text-sm font-medium'>6</span>
                  </div>
                  <div className='flex items-center justify-between'>
                    <div className='flex items-center gap-2'>
                      <div className='h-3 w-3 rounded-full bg-gray-500' />
                      <span className='text-sm'>
                        {t('dashboard.deliveredThisWeek')}
                      </span>
                    </div>
                    <span className='text-sm font-medium'>15</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Quick Actions */}
          <Card>
            <CardHeader>
              <CardTitle>{t('dashboard.quickActions')}</CardTitle>
            </CardHeader>
            <CardContent>
              <div className='grid grid-cols-2 gap-4 sm:grid-cols-4'>
                <a
                  href='/orders'
                  className='hover:bg-accent flex flex-col items-center gap-2 rounded-lg border p-4 transition-colors'
                >
                  <ClipboardList className='h-8 w-8' />
                  <span className='text-sm font-medium'>
                    {t('dashboard.newOrder')}
                  </span>
                </a>
                <a
                  href='/clients'
                  className='hover:bg-accent flex flex-col items-center gap-2 rounded-lg border p-4 transition-colors'
                >
                  <Users className='h-8 w-8' />
                  <span className='text-sm font-medium'>
                    {t('dashboard.addClient')}
                  </span>
                </a>
                <a
                  href='/prostheses'
                  className='hover:bg-accent flex flex-col items-center gap-2 rounded-lg border p-4 transition-colors'
                >
                  <Package className='h-8 w-8' />
                  <span className='text-sm font-medium'>
                    {t('dashboard.prosthesisCatalog')}
                  </span>
                </a>
                <a
                  href='/laboratories'
                  className='hover:bg-accent flex flex-col items-center gap-2 rounded-lg border p-4 transition-colors'
                >
                  <Building2 className='h-8 w-8' />
                  <span className='text-sm font-medium'>
                    {t('navigation.laboratories')}
                  </span>
                </a>
              </div>
            </CardContent>
          </Card>
        </div>
      </Main>
    </>
  )
}
