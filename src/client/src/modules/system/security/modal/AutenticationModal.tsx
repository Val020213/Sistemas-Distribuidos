'use client'
import { tailwindColors } from '@/theme/tailwindColors'
import { TabContext, TabPanel } from '@mui/lab'
import { Stack, Tab, Tabs } from '@mui/material'
import { useState } from 'react'
import Authentication from '../authentication/Authentication'
import HackerDialog from '@/core/dialog/HackerDIalog'
import Register from '../register/Register'

const AuthenticationModal = () => {
  const [selectedTab, setSelectedTab] = useState<'authentication' | 'register'>(
    'authentication'
  )
  return (
    <HackerDialog title={'PORTAL DE AUTENTICACIÓN'} onClose={() => {}} open>
      <Stack>
        <TabContext value={selectedTab}>
          <Stack px={3}>
            <Tabs
              value={selectedTab}
              onChange={(_, newValue) => setSelectedTab(newValue)}
              variant="fullWidth"
            >
              <Tab
                label="Autenticación"
                value={'authentication'}
                sx={{
                  fontFamily: 'monospace',
                  fontWeight: '400',
                  color:
                    selectedTab === 'authentication'
                      ? tailwindColors.green[500]
                      : tailwindColors.gray[500],
                }}
              />
              <Tab
                label="Registro"
                value={'register'}
                sx={{
                  fontFamily: 'monospace',
                  fontWeight: '400',
                  color:
                    selectedTab === 'register'
                      ? tailwindColors.green[500]
                      : tailwindColors.gray[500],
                }}
              />
            </Tabs>
          </Stack>
          <Stack>
            <TabPanel value={'authentication'}>
              <Authentication />
            </TabPanel>
            <TabPanel value={'register'}>
              <Register />
            </TabPanel>
          </Stack>
        </TabContext>
      </Stack>
    </HackerDialog>
  )
}

export default AuthenticationModal
