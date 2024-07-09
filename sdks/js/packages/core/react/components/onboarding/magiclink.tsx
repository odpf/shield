import { yupResolver } from '@hookform/resolvers/yup';
import { Button, Flex, Separator, Text, TextField } from '@raystack/apsara';
import React, { useCallback, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import * as yup from 'yup';
import { useFrontier } from '~/react/contexts/FrontierContext';

const styles = {
  container: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: 'var(--pd-16)'
  },

  button: {
    width: '100%'
  }
};

type MagicLinkProps = {
  open?: boolean;
  children?: React.ReactNode;
};

const emailSchema = yup.object({
  email: yup.string().trim().email().required()
});

type FormData = yup.InferType<typeof emailSchema>;

export const MagicLink = ({
  children,
  open = false,
  ...props
}: MagicLinkProps) => {
  const { client, config } = useFrontier();
  const [visiable, setVisiable] = useState<boolean>(open);
  const [loading, setLoading] = useState<boolean>(false);

  const {
    watch,
    control,
    handleSubmit,
    formState: { errors }
  } = useForm({
    resolver: yupResolver(emailSchema)
  });

  const magicLinkHandler = useCallback(
    async (data: FormData) => {
      setLoading(true);
      try {
        if (!client) return;

        const {
          data: { state = '' }
        } = await client.frontierServiceAuthenticate('mailotp', {
          email: data.email,
          callback_url: config.callbackUrl
        });

        const searchParams = new URLSearchParams({ state, email: data.email });

        // @ts-ignore
        window.location = `${
          config.redirectMagicLinkVerify
        }?${searchParams.toString()}`;
      } finally {
        setLoading(false);
      }
    },
    [client, config.callbackUrl, config.redirectMagicLinkVerify]
  );

  const email = watch('email', '');

  if (!visiable)
    return (
      <Button
        variant="secondary"
        size="medium"
        style={styles.button}
        onClick={() => setVisiable(true)}
      >
        Continue with Email
      </Button>
    );

  return (
    <form
      style={{ ...styles.container, flexDirection: 'column' }}
      onSubmit={handleSubmit(magicLinkHandler)}
    >
      {!open && <Separator />}
      <Flex
        direction={'column'}
        align={'start'}
        style={{
          width: '100%',
          position: 'relative',
          marginBottom: 'var(--pd-16)'
        }}
      >
        <Controller
          render={({ field }) => (
            <TextField
              {...field}
              // @ts-ignore
              size="medium"
              placeholder="name@example.com"
            />
          )}
          control={control}
          name="email"
        />

        <Text
          size={1}
          style={{
            color: 'var(--foreground-danger)',
            position: 'absolute',
            top: 'calc(100% + 4px)'
          }}
        >
          {errors.email && String(errors.email?.message)}
        </Text>
      </Flex>
      <Button
        size="medium"
        variant="primary"
        {...props}
        style={{ ...styles.button }}
        disabled={!email}
        type="submit"
      >
        {loading ? 'loading...' : 'Continue with Email'}
      </Button>
    </form>
  );
};
