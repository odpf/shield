'use client';

import { Button, Flex, Separator, Text } from '@raystack/apsara';
import { Outlet, useNavigate } from 'react-router-dom';
import { V1Beta1Organization } from '~/src';
import { styles } from '../styles';
import { GeneralProfile } from './general.profile';
import { GeneralOrganization } from './general.workspace';

type GeneralSettingProps = {
  organization?: V1Beta1Organization;
};

export default function GeneralSetting({ organization }: GeneralSettingProps) {
  return (
    <Flex direction="column" gap="large" style={{ width: '100%' }}>
      <Flex style={styles.header}>
        <Text size={6}>General</Text>
      </Flex>
      <Flex direction="column" gap="large" style={styles.container}>
        <GeneralProfile organization={organization} />
        <Separator />
        <GeneralOrganization organization={organization} />
        <Separator />
        <GeneralDeleteOrganization />
        <Separator />
      </Flex>
    </Flex>
  );
}

export const GeneralDeleteOrganization = () => {
  const navigate = useNavigate();
  return (
    <Flex direction="column" gap="medium">
      <Text size={3} style={{ color: 'var(--foreground-muted)' }}>
        If you want to permanently delete this organization and all of its data.
      </Text>

      <Button
        variant="danger"
        type="submit"
        size="medium"
        onClick={() => navigate('delete')}
      >
        Delete organization
      </Button>
      <Outlet />
    </Flex>
  );
};
