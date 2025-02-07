import { useMemo } from 'react'
import type { Hex } from 'viem'
import { useReadContract } from 'wagmi'
import { inbox, outbox } from '../index.js'
import type { OrderStatus } from '../types/orderStatus.js'

type UseOrderStatusParams = {
  destChainId: number
  originChainId?: number
  orderId?: Hex
  originData?: Hex
}

type UseDidFillParams = {
  destChainId: number
  orderId?: Hex
  originData?: Hex
}

function useDidFill(params: UseDidFillParams) {
  const { orderId, originData, destChainId } = params
  const filled = useReadContract({
    chainId: destChainId,
    address: outbox.address,
    abi: outbox.abi,
    functionName: 'didFill',
    args: orderId && originData ? [orderId, originData] : undefined,
    query: {
      enabled: !!orderId && !!originData,
      refetchInterval: 1000,
    },
  })

  return filled.data
}

export function useOrderStatus(params: UseOrderStatusParams) {
  const { originChainId, orderId } = params
  const filled = useDidFill({
    ...params,
  })

  const order = useReadContract({
    address: inbox.address,
    abi: inbox.abi,
    functionName: 'getOrder',
    chainId: originChainId,
    args: orderId ? [orderId] : undefined,
    query: {
      enabled: !!orderId || !filled,
      refetchInterval: 1000,
    },
  })

  const status: OrderStatus | undefined = useMemo(() => {
    return (
      order?.data &&
      (filled ? 'filled' : parseOrderStatus(order.data[1].status))
    )
  }, [order, filled])

  return status
}

const ORDER_STATUS: Record<number, OrderStatus> = {
  0: 'invalid',
  1: 'pending',
  2: 'accepted',
  3: 'rejected',
  4: 'reverted',
  5: 'filled',
  6: 'claimed',
} as const

function parseOrderStatus(status: number): OrderStatus {
  const orderStatus = ORDER_STATUS[status]
  if (!orderStatus) {
    throw new Error(`Invalid order status: ${status}`)
  }
  return orderStatus
}
