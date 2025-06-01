using System.Security.Cryptography;
using System.Text;

namespace Minitwit.Web.Helpers;

public static class Gravatar
{
    public static string GetGravatarUrl(string email, int size = 48)
    {
        email = email.Trim().ToLower();
        
        // Create an MD5 hash of the email
        using (var md5 = MD5.Create())
        {
            byte[] emailBytes = Encoding.UTF8.GetBytes(email);
            byte[] hashBytes = md5.ComputeHash(emailBytes);
            var hashString = BitConverter.ToString(hashBytes).Replace("-", "").ToLower();

            return string.Format("http://www.gravatar.com/avatar/{0}?d=identicon&s={1}", hashString, size);
        }
    }
}