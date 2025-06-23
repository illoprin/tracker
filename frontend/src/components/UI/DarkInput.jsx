function DarkInput({ children, ...props }) {
  return (
      <div className="relative">
        {children}

        <input
          {...props}
          className="w-[100%] h-[100%] bg-background-light text-text text-md pl-12 pt-3 pb-3 rounded-md"
        />
      </div>
  );
}

export default DarkInput;
