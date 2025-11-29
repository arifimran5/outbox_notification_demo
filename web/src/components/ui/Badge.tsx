import { cn } from '../../utils';

export const Badge = ({ children, variant = 'default', className }: { children: React.ReactNode, variant?: 'default' | 'secondary', className?: string }) => (
  <span className={cn(
    "inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
    variant === 'default' ? "border-transparent bg-blue-600 text-white shadow hover:bg-blue-700" : "border-transparent bg-gray-100 text-gray-900 hover:bg-gray-200",
    className
  )}>
    {children}
  </span>
);