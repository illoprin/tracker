import BackChevron from "@/components/shared/icons/BackChevron";

function BackButton({ children, ...props }) {
  return (
    <button
      className="text-secondary absolute top-5 left-5 hover:text-white cursor-pointer stroke-secondary hover:stroke-text flex items-center"
      {...props}
    >
      <BackChevron className="stroke-inherit" />
      <span className="block">
        {children}
      </span>
    </button>
  );
}

export default BackButton;
