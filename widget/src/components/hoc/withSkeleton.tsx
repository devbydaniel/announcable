interface WithLoadingProps {
  isLoading: boolean;
  skeleton: React.ReactNode;
}

export default function withSkeleton<T extends object>(
  WrappedComponent: React.ComponentType<T>,
) {
  return (props: T & WithLoadingProps) => {
    const { isLoading, skeleton, ...otherProps } = props;
    if (isLoading) {
      return skeleton;
    }
    return <WrappedComponent {...(otherProps as T)} />;
  };
}
