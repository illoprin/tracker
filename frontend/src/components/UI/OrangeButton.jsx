function OrangeButton({ children, ...props }) {
  return (
    <button className="btn-primary rounded-sm text-md p-2 pl-8 pr-8" {...props}>
      {children}
    </button>
  );
}

export default OrangeButton;
