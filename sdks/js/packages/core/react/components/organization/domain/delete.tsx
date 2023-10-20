import {
  Button,
  Dialog,
  Flex,
  Image,
  InputField,
  Separator,
  Text,
  TextField
} from '@raystack/apsara';

import { yupResolver } from '@hookform/resolvers/yup';
import { useNavigate, useParams } from '@tanstack/react-router';
import { useCallback, useEffect, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { toast } from 'sonner';
import * as yup from 'yup';
import cross from '~/react/assets/cross.svg';
import { useFrontier } from '~/react/contexts/FrontierContext';
import { V1Beta1Domain, V1Beta1Group } from '~/src';
import styles from '../organization.module.css';

const domainSchema = yup
  .object({
    domain: yup
      .string()
      .required()
      .matches(/[-a-zA-Z0-9.]{1,256}\.[a-zA-Z0-9()]{1,6}$/, 'Domain is invalid')
  })
  .required();

type FormData = yup.InferType<typeof domainSchema>;

export const DeleteDomain = () => {
  const {
    watch,
    control,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting }
  } = useForm({
    resolver: yupResolver(domainSchema)
  });
  const navigate = useNavigate({ from: '/domains/$domainId/delete' });
  const { domainId } = useParams({ from: '/domains/$domainId/delete' });
  const { client, activeOrganization: organization } = useFrontier();
  const [domain, setDomain] = useState<V1Beta1Domain>();
  const [isLoading, setIsLoading] = useState(false);

  const fetchDomainDetails = useCallback(async () => {
    if (!domainId) return;
    if (!organization?.id) return;

    try {
      setIsLoading(true);
      const {
        // @ts-ignore
        data: { domain }
      } = await client?.frontierServiceGetOrganizationDomain(
        organization?.id,
        domainId
      );
      setDomain(domain);
    } catch ({ error }: any) {
      toast.error('Something went wrong', {
        description: error.message
      });
    } finally {
      setIsLoading(false);
    }
  }, [client, domainId, organization?.id]);

  useEffect(() => {
    fetchDomainDetails();
  }, [fetchDomainDetails]);

  async function onSubmit(data: FormData) {
    // @ts-ignore. TODO: fix buf openapi plugin
    if (!domain?.id || !domain?.org_id) return;

    if (data.domain !== domain.name) {
      return setError('domain', { message: 'domain name is not same' });
    }
    try {
      await client?.frontierServiceDeleteOrganizationDomain(
        // @ts-ignore
        domain.org_id,
        domain.id
      );
      navigate({ to: '/domains' });
      toast.success('Domain deleted');
    } catch ({ error }: any) {
      toast.error('Something went wrong', {
        description: error.message
      });
    }
  }

  const domainName = watch('domain', '');

  return (
    <Dialog open={true}>
      <Dialog.Content
        style={{ padding: 0, maxWidth: '600px', width: '100%', zIndex: '60' }}
        overlayClassname={styles.overlay}
      >
        <Flex justify="between" style={{ padding: '16px 24px' }}>
          <Text size={6} style={{ fontWeight: '500' }}>
            Verify domain deletion
          </Text>
          <Image
            alt="cross"
            // @ts-ignore
            src={cross}
            onClick={() =>
              navigate({
                to: `/domains`
              })
            }
            style={{ cursor: 'pointer' }}
          />
        </Flex>
        <Separator />
        <form onSubmit={handleSubmit(onSubmit)}>
          <Flex
            direction="column"
            gap="medium"
            style={{ padding: '24px 32px' }}
          >
            <Text size={2}>
              This action can not be undone. This will permanently delete{' '}
              {domain?.name}
            </Text>

            <InputField label="Please type the domain name">
              <Controller
                render={({ field }) => (
                  <TextField
                    {...field}
                    // @ts-ignore
                    size="medium"
                    placeholder="Provide domain name"
                  />
                )}
                control={control}
                name="domain"
              />

              <Text size={1} style={{ color: 'var(--foreground-danger)' }}>
                {errors.domain && String(errors.domain?.message)}
              </Text>
            </InputField>
            <Flex>
              <Text size={2}>
                I acknowledge I understand that all of the team data will be
                deleted and want to proceed.
              </Text>
            </Flex>
            <Button
              variant="danger"
              size="medium"
              disabled={!domainName}
              type="submit"
              style={{ width: '100%' }}
            >
              {isSubmitting ? 'deleting...' : 'Delete this domain'}
            </Button>
          </Flex>
        </form>
      </Dialog.Content>
    </Dialog>
  );
};
