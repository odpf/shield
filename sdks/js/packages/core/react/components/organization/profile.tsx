import { Flex, ThemeProvider } from '@raystack/apsara';
import { useCallback, useEffect, useState } from 'react';
import {
  Outlet,
  RouterProvider,
  Router,
  Route,
  createMemoryHistory,
  useRouterContext,
  RouterContext
} from '@tanstack/react-router';
import { Toaster } from 'sonner';
import { useFrontier } from '~/react/contexts/FrontierContext';
import Domain from './domain';
import { AddDomain } from './domain/add-domain';
import { VerifyDomain } from './domain/verify-domain';
import GeneralSetting from './general';
import { DeleteOrganization } from './general/delete';
import WorkspaceMembers from './members';
import { InviteMember } from './members/invite';
import UserPreferences from './preferences';

import { default as WorkspaceProjects } from './project';
import { AddProject } from './project/add';
import { DeleteProject } from './project/delete';
import { ProjectPage } from './project/project';

import WorkspaceSecurity from './security';
import { Sidebar } from './sidebar';
import WorkspaceTeams from './teams';
import { AddTeam } from './teams/add';
import { DeleteTeam } from './teams/delete';
import { TeamPage } from './teams/team';
import { UserSetting } from './user';
import { SkeletonTheme } from 'react-loading-skeleton';
import { InviteTeamMembers } from './teams/members/invite';
import { DeleteDomain } from './domain/delete';
import Billing from './billing';
import Tokens from './tokens';
import { EditBillingAddress } from './billing/address/edit';
import { ConfirmCycleSwitch } from './billing/cycle-switch';
import Plans from './plans';
import ConfirmPlanChange from './plans/confirm-change';

interface OrganizationProfileProps {
  organizationId: string;
  defaultRoute?: string;
  showBilling?: boolean;
  showTokens?: boolean;
  showPreferences?: boolean;
  hideToast?: boolean;
}

const routerContext = new RouterContext<
  Pick<
    OrganizationProfileProps,
    | 'organizationId'
    | 'showBilling'
    | 'showTokens'
    | 'hideToast'
    | 'showPreferences'
  >
>();

const RootRouter = () => {
  const { organizationId, hideToast } = useRouterContext({ from: '__root__' });
  const {
    client,
    setActiveOrganization,
    setIsActiveOrganizationLoading,
    config
  } = useFrontier();

  const fetchOrganization = useCallback(async () => {
    try {
      setIsActiveOrganizationLoading(true);
      const {
        // @ts-ignore
        data: { organization }
      } = await client?.frontierServiceGetOrganization(organizationId);
      setActiveOrganization(organization);
    } catch (err) {
      console.error(err);
    } finally {
      setIsActiveOrganizationLoading(false);
    }
  }, [
    client,
    organizationId,
    setActiveOrganization,
    setIsActiveOrganizationLoading
  ]);

  useEffect(() => {
    if (organizationId) {
      fetchOrganization();
    } else {
      setActiveOrganization(undefined);
    }
  }, [organizationId, fetchOrganization, setActiveOrganization]);

  const visibleToasts = hideToast ? 0 : 1;

  return (
    <ThemeProvider defaultTheme={config?.theme}>
      <SkeletonTheme
        highlightColor="var(--background-base)"
        baseColor="var(--background-base-hover)"
      >
        <Toaster richColors visibleToasts={visibleToasts} />
        <Flex style={{ width: '100%', height: '100%' }}>
          <Sidebar />
          <Outlet />
        </Flex>
      </SkeletonTheme>
    </ThemeProvider>
  );
};

const rootRoute = routerContext.createRootRoute({
  component: RootRouter
});
const indexRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/',
  component: GeneralSetting
});

const deleteOrgRoute = new Route({
  getParentRoute: () => indexRoute,
  path: '/delete',
  component: DeleteOrganization
});

const securityRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/security',
  component: WorkspaceSecurity
});

const membersRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/members',
  component: WorkspaceMembers
});

const inviteMemberRoute = new Route({
  getParentRoute: () => membersRoute,
  path: '/modal',
  component: InviteMember
});

const teamsRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/teams',
  component: WorkspaceTeams
});

const addTeamRoute = new Route({
  getParentRoute: () => teamsRoute,
  path: '/modal',
  component: AddTeam
});

const domainsRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/domains',
  component: Domain
});

const verifyDomainRoute = new Route({
  getParentRoute: () => domainsRoute,
  path: '/$domainId/verify',
  component: VerifyDomain
});

const deleteDomainRoute = new Route({
  getParentRoute: () => domainsRoute,
  path: '/$domainId/delete',
  component: DeleteDomain
});

const addDomainRoute = new Route({
  getParentRoute: () => domainsRoute,
  path: '/modal',
  component: AddDomain
});

const teamRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/teams/$teamId',
  component: TeamPage
});

const inviteTeamMembersRoute = new Route({
  getParentRoute: () => teamRoute,
  path: '/invite',
  component: InviteTeamMembers
});

const deleteTeamRoute = new Route({
  getParentRoute: () => teamRoute,
  path: '/delete',
  component: DeleteTeam
});

const projectsRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/projects',
  component: WorkspaceProjects
});

const addProjectRoute = new Route({
  getParentRoute: () => projectsRoute,
  path: '/modal',
  component: AddProject
});

const projectPageRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/projects/$projectId',
  component: ProjectPage
});

const deleteProjectRoute = new Route({
  getParentRoute: () => projectPageRoute,
  path: '/delete',
  component: DeleteProject
});

const profileRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/profile',
  component: UserSetting
});

const preferencesRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/preferences',
  component: UserPreferences
});

const billingRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/billing',
  component: Billing
});

const editBillingAddressRoute = new Route({
  getParentRoute: () => billingRoute,
  path: '/$billingId/edit-address',
  component: EditBillingAddress
});

const switchBillingCycleModalRoute = new Route({
  getParentRoute: () => billingRoute,
  path: '/cycle-switch/$planId',
  component: ConfirmCycleSwitch
});

const plansRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/plans',
  component: Plans
});

const planDowngradeRoute = new Route({
  getParentRoute: () => plansRoute,
  path: '/confirm-change/$planId',
  component: ConfirmPlanChange
});

const tokensRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/tokens',
  component: Tokens
});

const routeTree = rootRoute.addChildren([
  indexRoute.addChildren([deleteOrgRoute]),
  securityRoute,
  membersRoute.addChildren([inviteMemberRoute]),
  teamsRoute.addChildren([addTeamRoute]),
  domainsRoute.addChildren([
    addDomainRoute,
    verifyDomainRoute,
    deleteDomainRoute
  ]),
  teamRoute.addChildren([deleteTeamRoute, inviteTeamMembersRoute]),
  projectsRoute.addChildren([addProjectRoute]),
  projectPageRoute.addChildren([deleteProjectRoute]),
  profileRoute,
  preferencesRoute,
  billingRoute.addChildren([
    editBillingAddressRoute,
    switchBillingCycleModalRoute
  ]),
  plansRoute.addChildren([planDowngradeRoute]),
  tokensRoute
]);

const router = new Router({
  routeTree,
  context: {
    organizationId: '',
    showBilling: false,
    showTokens: false,
    showPreferences: false
  }
});

export const OrganizationProfile = ({
  organizationId,
  defaultRoute = '/',
  showBilling = false,
  showTokens = false,
  showPreferences = false,
  hideToast = false
}: OrganizationProfileProps) => {
  const memoryHistory = createMemoryHistory({
    initialEntries: [defaultRoute]
  });

  const memoryRouter = new Router({
    routeTree,
    history: memoryHistory,
    context: {
      organizationId,
      showBilling,
      showTokens,
      hideToast,
      showPreferences
    }
  });
  return <RouterProvider router={memoryRouter} />;
};

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
