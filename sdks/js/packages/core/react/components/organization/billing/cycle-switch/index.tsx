import { Dialog, Flex, Text, Image, Separator, Button } from '@raystack/apsara';
import Skeleton from 'react-loading-skeleton';
import { useNavigate, useParams } from '@tanstack/react-router';
import cross from '~/react/assets/cross.svg';
import styles from '../../organization.module.css';
import { useCallback, useEffect, useState } from 'react';
import { useFrontier } from '~/react/contexts/FrontierContext';
import { V1Beta1Plan } from '~/src';
import { toast } from 'sonner';
import { getPlanIntervalName, getPlanPrice } from '~/react/utils';
import * as _ from 'lodash';
import { usePlans } from '../../plans/hooks/usePlans';
import dayjs from 'dayjs';
import { DEFAULT_DATE_FORMAT } from '~/react/utils/constants';

export function ConfirmCycleSwitch() {
  const { activePlan, client, paymentMethod, config } = useFrontier();
  const navigate = useNavigate({ from: '/billing/cycle-switch/$planId' });
  const { planId } = useParams({ from: '/billing/cycle-switch/$planId' });
  const dateFormat = config?.dateFormat || DEFAULT_DATE_FORMAT;

  const [isPlanLoading, setIsPlanLoading] = useState(false);
  const [nextPlan, setNextPlan] = useState<V1Beta1Plan>();
  const [isCycleSwitching, setCycleSwitching] = useState(false);

  const closeModal = useCallback(
    () => navigate({ to: '/billing' }),
    [navigate]
  );

  const {
    checkoutPlan,
    isLoading: isPlanActionLoading,
    changePlan,
    verifyPlanChange
  } = usePlans();

  const nextPlanPrice = nextPlan ? getPlanPrice(nextPlan) : { amount: 0 };
  const isPaymentMethodRequired =
    _.isEmpty(paymentMethod) && nextPlanPrice.amount > 0;

  const nextPlanIntervalName = getPlanIntervalName(nextPlan);

  const nextPlanMetadata = nextPlan?.metadata as Record<string, number>;
  const activePlanMetadata = activePlan?.metadata as Record<string, number>;

  const isUpgrade =
    (Number(nextPlanMetadata?.weightage) || 0) -
      (Number(activePlanMetadata?.weightage) || 0) >
    0;

  useEffect(() => {
    async function getNextPlan(nextPlanId: string) {
      setIsPlanLoading(true);
      try {
        const resp = await client?.frontierServiceGetPlan(nextPlanId);
        const plan = resp?.data?.plan;
        setNextPlan(plan);
      } catch (err: any) {
        toast.error('Something went wrong', {
          description: err.message
        });
        console.error(err);
      } finally {
        setIsPlanLoading(false);
      }
    }
    if (planId) {
      getNextPlan(planId);
    }
  }, [client, planId]);

  const isLoading = isPlanLoading;

  async function onConfirm() {
    setCycleSwitching(true);
    try {
      if (nextPlan?.id) {
        const nextPlanId = nextPlan?.id;
        if (isPaymentMethodRequired) {
          checkoutPlan({
            planId: nextPlanId,
            isTrial: false,
            onSuccess: data => {
              window.location.href = data?.checkout_url as string;
            }
          });
        } else
          changePlan({
            planId: nextPlanId,
            onSuccess: async () => {
              const planPhase = await verifyPlanChange({
                planId: nextPlanId
              });
              if (planPhase) {
                closeModal();
                const changeDate = dayjs(planPhase?.effective_at).format(
                  dateFormat
                );
                toast.success(`Plan cycle switch successful`, {
                  description: `Your plan cycle will switched to ${nextPlanIntervalName} on ${changeDate}`
                });
              }
            },
            immediate: isUpgrade
          });
      }
    } catch (err: any) {
      console.error(err);
      toast.error('Something went wrong', {
        description: err.message
      });
    } finally {
      setCycleSwitching(false);
    }
  }

  return (
    <Dialog open={true}>
      {/* @ts-ignore */}
      <Dialog.Content
        style={{ padding: 0, maxWidth: '600px', width: '100%', zIndex: '60' }}
        overlayClassname={styles.overlay}
      >
        <Flex justify="between" style={{ padding: '16px 24px' }}>
          <Text size={6} style={{ fontWeight: '500' }}>
            Switch billing cycle
          </Text>

          <Image
            alt="cross"
            style={{ cursor: 'pointer' }}
            // @ts-ignore
            src={cross}
            onClick={closeModal}
          />
        </Flex>
        <Separator />
        <Flex
          style={{ padding: 'var(--pd-32) 24px', gap: '24px' }}
          direction={'column'}
        >
          {isLoading ? (
            <Skeleton />
          ) : (
            <Flex gap="small">
              <Text size={2} weight={500}>
                Current cycle:
              </Text>
              <Text size={2} style={{ color: 'var(--foreground-muted)' }}>
                {getPlanIntervalName(activePlan)}
              </Text>
            </Flex>
          )}
          {isLoading ? (
            <Skeleton />
          ) : (
            <Flex gap="small">
              <Text size={2} weight={500}>
                New cycle:
              </Text>
              <Text size={2} style={{ color: 'var(--foreground-muted)' }}>
                {nextPlanIntervalName} (
                {isUpgrade
                  ? 'effective immediately'
                  : 'effective from the next billing cycle'}
                )
              </Text>
            </Flex>
          )}
        </Flex>
        <Separator />
        <Flex justify={'end'} gap="medium" style={{ padding: 'var(--pd-16)' }}>
          <Button variant={'secondary'} onClick={closeModal} size={'medium'}>
            Cancel
          </Button>
          <Button
            variant={'primary'}
            size={'medium'}
            disabled={isLoading || isCycleSwitching || isPlanActionLoading}
            onClick={onConfirm}
          >
            {isCycleSwitching ? 'Switching...' : 'Switch cycle'}
          </Button>
        </Flex>
      </Dialog.Content>
    </Dialog>
  );
}
