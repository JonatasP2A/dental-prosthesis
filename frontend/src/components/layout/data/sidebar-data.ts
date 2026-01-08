import { useTranslation } from 'react-i18next'
import {
  LayoutDashboard,
  Monitor,
  Bell,
  Palette,
  Settings,
  Wrench,
  UserCog,
  Users,
  Building2,
  ClipboardList,
  Package,
  HardHat,
} from 'lucide-react'
import { type SidebarData } from '../types'

export function useSidebarData(): SidebarData {
  const { t } = useTranslation()

  return {
    user: {
      name: 'User',
      email: 'user@example.com',
      avatar: '/avatars/shadcn.jpg',
    },
    teams: [
      {
        name: t('sidebar.team.dentalLab'),
        logo: Building2,
        plan: t('sidebar.team.professional'),
      },
    ],
    navGroups: [
      {
        title: t('navigation.general'),
        items: [
          {
            title: t('navigation.dashboard'),
            url: '/',
            icon: LayoutDashboard,
          },
          {
            title: t('navigation.laboratories'),
            url: '/laboratories',
            icon: Building2,
          },
          {
            title: t('navigation.clients'),
            url: '/clients',
            icon: Users,
          },
          {
            title: t('navigation.orders'),
            url: '/orders',
            icon: ClipboardList,
          },
          {
            title: t('navigation.prostheses'),
            url: '/prostheses',
            icon: Package,
          },
          {
            title: t('navigation.technicians'),
            url: '/technicians',
            icon: HardHat,
          },
        ],
      },
      {
        title: t('navigation.settings'),
        items: [
          {
            title: t('navigation.settings'),
            icon: Settings,
            items: [
              {
                title: t('navigation.profile'),
                url: '/settings',
                icon: UserCog,
              },
              {
                title: t('navigation.account'),
                url: '/settings/account',
                icon: Wrench,
              },
              {
                title: t('navigation.appearance'),
                url: '/settings/appearance',
                icon: Palette,
              },
              {
                title: t('navigation.notifications'),
                url: '/settings/notifications',
                icon: Bell,
              },
              {
                title: t('navigation.display'),
                url: '/settings/display',
                icon: Monitor,
              },
            ],
          },
        ],
      },
    ],
  }
}

// Keep the static export for backward compatibility but it won't have translations
export const sidebarData: SidebarData = {
  user: {
    name: 'User',
    email: 'user@example.com',
    avatar: '/avatars/shadcn.jpg',
  },
  teams: [
    {
      name: 'Dental Lab',
      logo: Building2,
      plan: 'Professional',
    },
  ],
  navGroups: [
    {
      title: 'General',
      items: [
        {
          title: 'Dashboard',
          url: '/',
          icon: LayoutDashboard,
        },
        {
          title: 'Laboratories',
          url: '/laboratories',
          icon: Building2,
        },
        {
          title: 'Clients',
          url: '/clients',
          icon: Users,
        },
        {
          title: 'Orders',
          url: '/orders',
          icon: ClipboardList,
        },
        {
          title: 'Prostheses',
          url: '/prostheses',
          icon: Package,
        },
        {
          title: 'Technicians',
          url: '/technicians',
          icon: HardHat,
        },
      ],
    },
    {
      title: 'Settings',
      items: [
        {
          title: 'Settings',
          icon: Settings,
          items: [
            {
              title: 'Profile',
              url: '/settings',
              icon: UserCog,
            },
            {
              title: 'Account',
              url: '/settings/account',
              icon: Wrench,
            },
            {
              title: 'Appearance',
              url: '/settings/appearance',
              icon: Palette,
            },
            {
              title: 'Notifications',
              url: '/settings/notifications',
              icon: Bell,
            },
            {
              title: 'Display',
              url: '/settings/display',
              icon: Monitor,
            },
          ],
        },
      ],
    },
  ],
}
