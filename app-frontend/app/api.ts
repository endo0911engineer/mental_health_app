const BASE_URL = 'http://localhost:8080';

// AIモデルへのリクエスト
export const analyzeSentiment = async (sentence: string) => {

      const response = await fetch(`${BASE_URL}/predict`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ sentence }),
      });

      if (!response.ok) {
        throw new Error("Failed to fetch the analysis result");
      }

      const data = await response.json();
      console.log("Response from Go server:", data);
  
      return data;
}

// サインアップ
export const signUp = async (userData: { email: string; password: string }) => {
  const response = await fetch(`${BASE_URL}/signup`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    throw new Error(await response.text());
  }

  return await response.json()
};

// ログイン
export const logIn = async (userData: { email: string; password: string }) => {
  const response = await fetch(`${BASE_URL}/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  // ステータスコードをチェックして詳細なエラーを確認
  if (!response.ok) {
    const errorText = await response.text();
    console.error('Error status:', response.status, 'Error text:', errorText);
    throw new Error(errorText || 'Unexpected error');
  }

  return await response.json()
};

// トークンをヘッダーに設定するための関数
export const setAuthToken = (token: string) => {
  return {
    'Authorization': `Bearer ${token}`,
  };
};


export const saveEmotion = async (date: string, emotionText: string, userId: number) => {
  const isoDate = new Date(date).toISOString();

  const response = await fetch(`${BASE_URL}/dashboard/saveEmotion?user_id=${userId}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
    body: JSON.stringify({
      date: isoDate,
      emotion: emotionText,
    }),
  });

  console.log(response);

  if (!response.ok) {
    throw new Error('Failed to save emotion');
  }

  return await response.json();
};

export const updateEmotion = async (date: string, emotionText: string, userId: number) => {
  const isoDate = new Date(date).toISOString();

  const response = await fetch(`${BASE_URL}/dashboard/updateEmotion?user_id=${userId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
    body: JSON.stringify({
      date: isoDate,
      emotion: emotionText,
    }),
  });

  console.log(response);

  if (!response.ok) {
    throw new Error('Failed to save emotion');
  }

  return await response.json();
};

export const getEmotions = async (userId: number) => {
  const response = await fetch(`${BASE_URL}/dashboard/getEmotions?user_id=${userId}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  });

  console.log(response)
  if (!response.ok) {
    throw new Error('Failed to get emotion');
  }

  return await response.json();
};

export const deleteEmotion= async (date: string, userId: number) => {
  const response = await fetch(`${BASE_URL}/dashboard/deleteEmotion?user_id=${userId}&date=${encodeURIComponent(date)}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  });

  const responseText = await response.text(); // レスポンスをテキストとして取得

  console.log(responseText)
  if (!response.ok) {
    throw new Error('Failed to delete emotion');
  }

  return responseText;
};

export const generateGraph = async (userId: number, emotionText: string) => {
  const response = await fetch(`${BASE_URL}/dashboard/generateGraph?user_id=${userId}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
    body: JSON.stringify({ emotion: emotionText})
  });

  console.log(response)
  if (!response.ok) {
    throw new Error('Error to get emotion');
  }

  return await response.json()
} 