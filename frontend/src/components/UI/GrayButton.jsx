function GrayButton({ children, ...props }) {
  return (
    <button
      className="bg-background-secondary hover:bg-background-light transition-colors rounded-sm text-md p-2"
      {...props}
    >
      {children}
    </button>
  );
}

export default GrayButton;
