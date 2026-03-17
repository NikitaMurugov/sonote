import { ref } from 'vue'

export function useDebounce<T extends (...args: any[]) => any>(fn: T, delay: number = 300) {
  let timeout: ReturnType<typeof setTimeout>

  const debounced = (...args: Parameters<T>) => {
    clearTimeout(timeout)
    timeout = setTimeout(() => fn(...args), delay)
  }

  return debounced
}

export function useDebouncedRef<T>(initialValue: T, delay: number = 300) {
  const value = ref(initialValue) as ReturnType<typeof ref<T>>
  let timeout: ReturnType<typeof setTimeout>

  const setValue = (newValue: T) => {
    clearTimeout(timeout)
    timeout = setTimeout(() => {
      value.value = newValue as any
    }, delay)
  }

  return { value, setValue }
}
