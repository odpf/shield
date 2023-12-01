import {
  DotsHorizontalIcon,
  TrashIcon,
  UpdateIcon
} from '@radix-ui/react-icons';
import { Avatar, DropdownMenu, Flex, Label, Text } from '@raystack/apsara';
import { useNavigate, useParams } from '@tanstack/react-router';
import type { ColumnDef } from '@tanstack/react-table';
import Skeleton from 'react-loading-skeleton';
import { toast } from 'sonner';
import { useFrontier } from '~/react/contexts/FrontierContext';
import { V1Beta1Role, V1Beta1User } from '~/src';
import { Role } from '~/src/types';
import { differenceWith, getInitials, isEqualById } from '~/utils';
import styles from '../../organization.module.css';

interface getColumnsOptions {
  roles: V1Beta1Role[];
  organizationId: string;
  canUpdateGroup?: boolean;
  memberRoles?: Record<string, Role[]>;
  isLoading?: boolean;
}

export const getColumns: (
  options: getColumnsOptions
) => ColumnDef<V1Beta1User, any>[] = ({
  roles = [],
  organizationId,
  canUpdateGroup = false,
  memberRoles = {},
  isLoading
}) => [
  {
    header: '',
    accessorKey: 'avatar',
    size: 44,
    meta: {
      style: {
        width: '30px',
        padding: 0
      }
    },
    cell: isLoading
      ? () => <Skeleton />
      : ({ row, getValue }) => {
          return (
            <Avatar
              src={getValue()}
              fallback={getInitials(row.original?.title || row.original?.email)}
              // @ts-ignore
              style={{ marginRight: 'var(--mr-12)' }}
            />
          );
        }
  },
  {
    header: 'Title',
    accessorKey: 'title',
    meta: {
      style: {
        paddingLeft: 0
      }
    },
    cell: isLoading
      ? () => <Skeleton />
      : ({ row, getValue }) => {
          return (
            <Flex direction="column" gap="extra-small">
              <Label style={{ fontWeight: '$500' }}>{getValue()}</Label>
              <Text>{row.original.email}</Text>
            </Flex>
          );
        }
  },
  {
    header: 'Roles',
    accessorKey: 'email',
    cell: isLoading
      ? () => <Skeleton />
      : ({ row, getValue }) => {
          return (
            (row.original?.id &&
              memberRoles[row.original?.id] &&
              memberRoles[row.original?.id]
                .map((r: any) => r.title || r.name)
                .join(', ')) ??
            'Inherited role'
          );
        }
  },
  {
    header: '',
    accessorKey: 'id',
    meta: {
      style: {
        textAlign: 'end'
      }
    },
    cell: isLoading
      ? () => <Skeleton />
      : ({ row }) => (
          <MembersActions
            member={row.original as V1Beta1User}
            organizationId={organizationId}
            canUpdateGroup={canUpdateGroup}
            excludedRoles={differenceWith<V1Beta1Role>(
              isEqualById,
              roles,
              row.original?.id && memberRoles[row.original?.id]
                ? memberRoles[row.original?.id]
                : []
            )}
          />
        )
  }
];

const MembersActions = ({
  member,
  organizationId,
  canUpdateGroup,
  excludedRoles = []
}: {
  member: V1Beta1User;
  canUpdateGroup?: boolean;
  organizationId: string;
  excludedRoles: V1Beta1Role[];
}) => {
  let { teamId } = useParams({ from: '/teams/$teamId' });
  const { client } = useFrontier();
  const navigate = useNavigate({ from: '/teams/$teamId' });

  async function deleteMember() {
    try {
      await client?.frontierServiceRemoveGroupUser(
        organizationId,
        teamId as string,
        member?.id as string
      );
      navigate({
        to: '/teams/$teamId',
        params: {
          teamId
        }
      });
      toast.success('Member deleted');
    } catch ({ error }: any) {
      toast.error('Something went wrong', {
        description: error.message
      });
    }
  }

  async function updateRole() {
    try {
      toast.success('Member role deleted');
    } catch ({ error }: any) {
      toast.error('Something went wrong', {
        description: error.message
      });
    }
  }
  return canUpdateGroup ? (
    <DropdownMenu>
      <DropdownMenu.Trigger asChild style={{ cursor: 'pointer' }}>
        <DotsHorizontalIcon />
      </DropdownMenu.Trigger>
      <DropdownMenu.Content align="end">
        <DropdownMenu.Group>
          {excludedRoles.map((role: V1Beta1Role) => (
            <DropdownMenu.Item style={{ padding: 0 }} key={role.id}>
              <div onClick={updateRole} className={styles.dropdownActionItem}>
                <UpdateIcon />
                Make {role.title}
              </div>
            </DropdownMenu.Item>
          ))}
          <DropdownMenu.Item style={{ padding: 0 }}>
            <div onClick={deleteMember} className={styles.dropdownActionItem}>
              <TrashIcon />
              Remove from team
            </div>
          </DropdownMenu.Item>
        </DropdownMenu.Group>
      </DropdownMenu.Content>
    </DropdownMenu>
  ) : null;
};
