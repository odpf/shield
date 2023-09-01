import { Flex, Text } from '@raystack/apsara';

import { Tabs } from '@raystack/apsara';
import { Outlet, useParams } from '@tanstack/react-router';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { useFrontier } from '~/react/contexts/FrontierContext';
import { V1Beta1Group, V1Beta1User } from '~/src';
import { styles } from '../styles';
import { General } from './general';
import { Members } from './members';

export const TeamPage = () => {
  let { teamId } = useParams({ from: '/teams/$teamId' });
  const [team, setTeam] = useState<V1Beta1Group>();
  const [orgMembers, setOrgMembers] = useState<V1Beta1User[]>([]);
  const [members, setMembers] = useState<V1Beta1User[]>([]);
  const { client, activeOrganization: organization } = useFrontier();

  useEffect(() => {
    async function getTeamDetails() {
      if (!organization?.id || !teamId) return;

      try {
        const {
          // @ts-ignore
          data: { group }
        } = await client?.frontierServiceGetGroup(organization?.id, teamId);

        setTeam(group);
      } catch ({ error }: any) {
        toast.error('Something went wrong', {
          description: error.message
        });
      }
    }
    getTeamDetails();
  }, [client, organization?.id, teamId]);

  useEffect(() => {
    async function getTeamMembers() {
      if (!organization?.id || !teamId) return;
      try {
        const {
          // @ts-ignore
          data: { users }
        } = await client?.frontierServiceListGroupUsers(
          organization?.id,
          teamId
        );
        setMembers(users);
      } catch ({ error }: any) {
        toast.error('Something went wrong', {
          description: error.message
        });
      }
    }
    getTeamMembers();
  }, [client, organization?.id, teamId]);

  useEffect(() => {
    async function getOrganizationMembers() {
      if (!organization?.id) return;
      try {
        const {
          // @ts-ignore
          data: { users }
        } = await client?.frontierServiceListOrganizationUsers(
          organization?.id
        );
        setOrgMembers(users);
      } catch ({ error }: any) {
        toast.error('Something went wrong', {
          description: error.message
        });
      }
    }
    getOrganizationMembers();
  }, [client, organization?.id]);

  return (
    <Flex direction="column" style={{ width: '100%' }}>
      <Flex style={styles.header}>
        <Text size={6}>Teams</Text>
      </Flex>
      <Tabs defaultValue="general" style={styles.container}>
        <Tabs.List elevated>
          <Tabs.Trigger value="general" style={{ flex: 1, height: 24 }}>
            General
          </Tabs.Trigger>
          <Tabs.Trigger value="members" style={{ flex: 1, height: 24 }}>
            Members
          </Tabs.Trigger>
        </Tabs.List>
        <Tabs.Content value="general">
          <General organization={organization} team={team} />
        </Tabs.Content>
        <Tabs.Content value="members">
          <Members
            orgMembers={orgMembers}
            members={members}
            setMembers={setMembers}
            organizationId={organization?.id}
          />
        </Tabs.Content>
      </Tabs>
      <Outlet />
    </Flex>
  );
};
