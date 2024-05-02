import qs from 'query-string';
import { useFrontier } from '~/react/contexts/FrontierContext';
import * as _ from 'lodash';
import { Button, Flex, Text } from '@raystack/apsara';
import Skeleton from 'react-loading-skeleton';
import billingStyles from './billing.module.css';
import { V1Beta1PaymentMethod } from '~/src';
import { toast } from 'sonner';
import { useState } from 'react';

interface PaymentMethodProps {
  paymentMethod?: V1Beta1PaymentMethod;
  isLoading: boolean;
}

export const PaymentMethod = ({
  paymentMethod = {},
  isLoading
}: PaymentMethodProps) => {
  const { client, config, billingAccount } = useFrontier();
  const [isActionLoading, setIsActionLoading] = useState(false);
  const {
    card_last4 = '',
    card_expiry_month,
    card_expiry_year
  } = paymentMethod;
  // TODO: change card digit as per card type
  const cardDigit = 12;
  const cardNumber = card_last4 ? _.repeat('*', cardDigit) + card_last4 : 'N/A';
  const cardExp =
    card_expiry_month && card_expiry_year
      ? `${card_expiry_month}/${card_expiry_year}`
      : 'N/A';

  const isPaymentMethodAvailable = card_last4 !== '';

  const updatePaymentMethod = () => {};

  const addPaymentMethod = async () => {
    const orgId = billingAccount?.org_id || '';
    const billingAccountId = billingAccount?.id || '';
    if (billingAccountId && orgId) {
      setIsActionLoading(true);
      try {
        const query = qs.stringify(
          {
            details: btoa(
              qs.stringify({
                billing_id: billingAccount?.id,
                organization_id: billingAccount?.org_id,
                type: 'billing'
              })
            ),
            checkout_id: '{{.CheckoutID}}'
          },
          { encode: false }
        );
        const cancel_url = `${config?.billing?.cancelUrl}?${query}`;
        const success_url = `${config?.billing?.successUrl}?${query}`;

        const resp = await client?.frontierServiceCreateCheckout(
          billingAccount?.org_id || '',
          billingAccount?.id || '',
          {
            cancel_url,
            success_url,
            setup_body: {
              payment_method: true
            }
          }
        );
        const checkout_url = resp?.data?.checkout_session?.checkout_url;
        if (checkout_url) {
          window.location.href = checkout_url;
        }
      } catch (err) {
        console.error(err);
        toast.error('Something went wrong');
      } finally {
        setIsActionLoading(true);
      }
    }
  };

  function onClick() {
    if (isPaymentMethodAvailable) {
      updatePaymentMethod();
    } else {
      addPaymentMethod;
    }
  }

  const isBtnDisabled = isLoading || isActionLoading;

  return (
    <div className={billingStyles.detailsBox}>
      <Flex align={'center'} justify={'between'} style={{ width: '100%' }}>
        <Text className={billingStyles.detailsBoxHeading}>Payment method</Text>
        <Button
          variant={'secondary'}
          onClick={onClick}
          disabled={isBtnDisabled}
        >
          {isPaymentMethodAvailable ? 'Update' : 'Add method'}
        </Button>
      </Flex>
      <Flex direction={'column'} gap={'extra-small'}>
        <Text className={billingStyles.detailsBoxRowLabel}>Card Number</Text>
        <Text className={billingStyles.detailsBoxRowValue}>
          {isLoading ? <Skeleton /> : cardNumber}
        </Text>
      </Flex>
      <Flex direction={'column'} gap={'extra-small'}>
        <Text className={billingStyles.detailsBoxRowLabel}>Expiry</Text>
        <Text className={billingStyles.detailsBoxRowValue}>
          {isLoading ? <Skeleton /> : cardExp}
        </Text>
      </Flex>
    </div>
  );
};
