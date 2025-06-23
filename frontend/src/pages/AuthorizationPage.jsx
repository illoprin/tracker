import LoginForm from '@/components/pages/authorization/LoginForm';
import RegForm from '@/components/pages/authorization/RegForm';
import EnvelopeIcon from '@/components/shared/icons/EnvelopeIcon';
import LockIcon from '@/components/shared/icons/LockIcon';
import TrackerLogo from '@/components/shared/icons/Logo';
import UserIcon from '@/components/shared/icons/UserIcon';
import BackButton from '@/components/UI/BackButton';
import DarkInput from '@/components/UI/DarkInput';
import OrangeButton from '@/components/UI/OrangeButton';
import { routes } from '@/router';
import { useEffect, useMemo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

function AuthorizationPage() {
  const { hash } = useLocation();
  const navigate = useNavigate();
  

  const mode = () => {
    if (hash == routes.public.auth.hashes.auth) {
      return 'auth';
    } else {
      return 'reg';
    }
  };

  return (
    <div className="flex items-center justify-center h-[100%]">
      <div className="bg-background-secondary overflow-hidden rounded-lg border-background-light flex panel border">
        {/* <!-- Logo Side --> */}
        <div className="relative panel-aside p-5 flex items-center justify-center">
          {/* <!-- Back Button --> */}
          <BackButton onClick={() => navigate(-1)}>Назад</BackButton>

          {/* <!-- Logo --> */}
          <TrackerLogo className="block ml-[4rem] mr-[4rem] w-[12rem]" />
        </div>

        {/* <!-- Form Side --> */}
        <div className="bg-background-secondary p-form max-w-[450px] min-h-[600px] flex flex-col justify-between">
          {mode() == 'reg' ? <RegForm /> : <LoginForm />}
        </div>
      </div>
    </div>
  );
}

export default AuthorizationPage;
