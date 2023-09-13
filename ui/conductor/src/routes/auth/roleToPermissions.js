export const roleToPermissions = {
  admin: {
    fullAccess: [
      '/auth',
      '/auth/users',
      '/auth/third-party',
      '/devices',
      '/devices/.*',
      '/ops',
      '/ops/overview',
      '/ops/emergency-lighting',
      '/ops/notifications',
      '/ops/security',
      '/site',
      '/site/zone',
      '/site/zone/.*',
      '/automations',
      '/automations/.*',
      '/system',
      '/system/drivers',
      '/system/features',
      '/system/components',
      '/signage'
    ],
    limitedAccess: [
    ],
    blockedAccess: [
    ]
  },
  commissioner: {
    fullAccess: [
      '/devices',
      '/devices/.*',
      '/ops',
      '/ops/overview',
      '/ops/emergency-lighting',
      '/ops/notifications',
      '/ops/security',
      '/site',
      '/site/zone',
      '/site/zone/.*',
      '/automations',
      '/automations/.*',
      '/system',
      '/system/drivers',
      '/system/features',
      '/system/components',
      '/signage'
    ],
    limitedAccess: [
    ],
    blockedAccess: [
      '/auth',
      '/auth/users',
      '/auth/third-party'
    ]
  },
  operator: {
    fullAccess: [
      '/devices',
      '/devices/.*',
      '/ops',
      '/ops/overview',
      '/ops/emergency-lighting',
      '/ops/notifications',
      '/ops/security',
      '/signage'
    ],
    limitedAccess: [
      '/site',
      '/site/zone',
      '/site/zone/.*',
      '/automations',
      '/automations/.*',
      '/system',
      '/system/drivers',
      '/system/features',
      '/system/components'
    ],
    blockedAccess: [
      '/auth',
      '/auth/users',
      '/auth/third-party'
    ]
  },
  signage: {
    fullAccess: [
      '/signage'
    ],
    limitedAccess: [
    ],
    blockedAccess: [
      '/auth',
      '/auth/users',
      '/auth/third-party',
      '/devices',
      '/devices/.*',
      '/ops',
      '/ops/overview',
      '/ops/emergency-lighting',
      '/ops/notifications',
      '/ops/security',
      '/site',
      '/site/zone',
      '/site/zone/.*',
      '/automations',
      '/automations/.*',
      '/system',
      '/system/drivers',
      '/system/features',
      '/system/components'
    ]
  },
  superAdmin: {
    fullAccess: [
      '/auth',
      '/auth/users',
      '/auth/third-party',
      '/devices',
      '/devices/.*',
      '/ops',
      '/ops/overview',
      '/ops/emergency-lighting',
      '/ops/notifications',
      '/ops/security',
      '/site',
      '/site/zone',
      '/site/zone/.*',
      '/automations',
      '/automations/.*',
      '/system',
      '/system/drivers',
      '/system/features',
      '/system/components',
      '/signage'
    ],
    limitedAccess: [
    ],
    blockedAccess: [
    ]
  },
  viewer: {
    fullAccess: [
      '/signage'
    ],
    limitedAccess: [
      '/auth',
      '/auth/users',
      '/auth/third-party',
      '/devices',
      '/devices/.*',
      '/ops',
      '/ops/overview',
      '/ops/emergency-lighting',
      '/ops/notifications',
      '/ops/security',
      '/site',
      '/site/zone',
      '/site/zone/.*',
      '/automations',
      '/automations/.*',
      '/system',
      '/system/drivers',
      '/system/features',
      '/system/components'
    ],
    blockedAccess: [
    ]
  }
};
