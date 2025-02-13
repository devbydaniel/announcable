import TanstackProvider from "./tanstack";

interface Props {
  children: React.ReactNode;
}

export default function Providers({ children }: Props) {
  return <TanstackProvider>{children}</TanstackProvider>;
}
